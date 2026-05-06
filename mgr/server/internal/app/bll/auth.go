package bll

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/dal"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"clearbill/mgr/server/pkg/passwordx"
	"gorm.io/gorm"
)

const sessionTTL = 24 * time.Hour

type AuthService struct {
	UserDAL       *dal.UserDAL
	SessionDAL    *dal.SessionDAL
	APITokenDAL   *dal.APITokenDAL
	CredentialDAL *dal.CredentialDAL
}

func NewAuthService(userDAL *dal.UserDAL, sessionDAL *dal.SessionDAL, apiTokenDAL *dal.APITokenDAL, credentialDAL *dal.CredentialDAL) *AuthService {
	return &AuthService{
		UserDAL:       userDAL,
		SessionDAL:    sessionDAL,
		APITokenDAL:   apiTokenDAL,
		CredentialDAL: credentialDAL,
	}
}

func (s *AuthService) Login(ctx context.Context, req *vo.LoginReq) (*dbmodel.User, string, error) {
	user, err := s.UserDAL.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, "", err
	}
	if user.Status != dbmodel.StatusActive {
		return nil, "", errors.New("user is disabled")
	}
	if err := passwordx.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	sessionToken, err := generateToken()
	if err != nil {
		return nil, "", err
	}

	session := &dbmodel.UserSession{
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(sessionTTL),
	}
	if err := s.SessionDAL.Create(ctx, session); err != nil {
		return nil, "", err
	}

	return user, sessionToken, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.SessionDAL.DeleteByToken(ctx, token)
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*dbmodel.User, error) {
	session, err := s.SessionDAL.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if session.ExpiresAt.Before(time.Now()) {
		_ = s.SessionDAL.DeleteByToken(ctx, token)
		return nil, gorm.ErrRecordNotFound
	}

	return s.UserDAL.GetByID(ctx, session.UserID)
}

func (s *AuthService) ValidateAPIToken(ctx context.Context, token string) (*dbmodel.User, error) {
	item, err := s.APITokenDAL.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return s.UserDAL.GetByID(ctx, item.UserID)
}

func (s *AuthService) ValidateCredentialToken(ctx context.Context, token string) (*dbmodel.User, error) {
	credential, err := s.CredentialDAL.GetActiveByToken(ctx, token)
	if err == nil {
		_ = s.CredentialDAL.TouchLastUsed(ctx, credential.ID)
		return s.UserDAL.GetByID(ctx, credential.UserID)
	}

	item, legacyErr := s.APITokenDAL.GetByToken(ctx, token)
	if legacyErr != nil {
		return nil, legacyErr
	}
	return s.UserDAL.GetByID(ctx, item.UserID)
}

func (s *AuthService) ValidateAccessKeySecret(ctx context.Context, accessKey, secretKey string) (*dbmodel.User, error) {
	credential, err := s.CredentialDAL.GetActiveByAccessKey(ctx, accessKey)
	if err != nil {
		return nil, err
	}
	if credential.Type != dbmodel.CredentialTypeAKSK || credential.SecretKey == nil {
		return nil, gorm.ErrRecordNotFound
	}
	if strings.TrimSpace(*credential.SecretKey) != strings.TrimSpace(secretKey) {
		return nil, gorm.ErrRecordNotFound
	}

	_ = s.CredentialDAL.TouchLastUsed(ctx, credential.ID)
	return s.UserDAL.GetByID(ctx, credential.UserID)
}

func (s *AuthService) ValidateAccessKeySignature(ctx context.Context, accessKey string, timestamp int64, signature string, method string, requestPath string, rawQuery string, body []byte) (*dbmodel.User, error) {
	credential, err := s.CredentialDAL.GetActiveByAccessKey(ctx, accessKey)
	if err != nil {
		return nil, err
	}
	if credential.Type != dbmodel.CredentialTypeAKSK || credential.SecretKey == nil {
		return nil, gorm.ErrRecordNotFound
	}
	if !isCredentialTimestampFresh(timestamp) {
		return nil, gorm.ErrRecordNotFound
	}

	expected := signCredentialRequest(*credential.SecretKey, timestamp, method, requestPath, rawQuery, body)
	if !hmac.Equal([]byte(expected), []byte(strings.ToLower(strings.TrimSpace(signature)))) {
		return nil, gorm.ErrRecordNotFound
	}

	_ = s.CredentialDAL.TouchLastUsed(ctx, credential.ID)
	return s.UserDAL.GetByID(ctx, credential.UserID)
}

func (s *AuthService) CreateAPIToken(ctx context.Context, userID uint, req *vo.CreateAPITokenReq) (*vo.CreateAPITokenResp, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	item := &dbmodel.UserAPIToken{
		UserID: userID,
		Name:   req.Name,
		Token:  token,
	}
	if err := s.APITokenDAL.Create(ctx, item); err != nil {
		return nil, err
	}

	return &vo.CreateAPITokenResp{
		Name:  item.Name,
		Token: item.Token,
	}, nil
}

func (s *AuthService) ChangeOwnPassword(ctx context.Context, userID uint, req *vo.ChangeOwnPasswordReq) error {
	user, err := s.UserDAL.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := passwordx.ComparePassword(user.PasswordHash, req.OldPassword); err != nil {
		return errors.New("old password is incorrect")
	}
	if err := passwordx.ComparePassword(user.PasswordHash, req.NewPassword); err == nil {
		return errors.New("new password must be different from old password")
	}
	hash, err := passwordx.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.UserDAL.Update(ctx, user)
}

func generateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func isCredentialTimestampFresh(ts int64) bool {
	if ts <= 0 {
		return false
	}
	now := time.Now().Unix()
	const skew = int64((5 * time.Minute) / time.Second)
	if ts < now-skew || ts > now+skew {
		return false
	}
	return true
}

func signCredentialRequest(secretKey string, timestamp int64, method, requestPath, rawQuery string, body []byte) string {
	payload := fmt.Sprintf("%s\n%s\n%s\n%d\n%s",
		strings.ToUpper(strings.TrimSpace(method)),
		strings.TrimSpace(requestPath),
		strings.TrimSpace(rawQuery),
		timestamp,
		bodyDigest(body),
	)
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, _ = io.WriteString(mac, payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func bodyDigest(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}
