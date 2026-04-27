package bll

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/dal"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"clearbill/mgr/server/pkg/passwordx"
	"gorm.io/gorm"
)

const sessionTTL = 24 * time.Hour

type AuthService struct {
	UserDAL     *dal.UserDAL
	SessionDAL  *dal.SessionDAL
	APITokenDAL *dal.APITokenDAL
}

func NewAuthService(userDAL *dal.UserDAL, sessionDAL *dal.SessionDAL, apiTokenDAL *dal.APITokenDAL) *AuthService {
	return &AuthService{
		UserDAL:     userDAL,
		SessionDAL:  sessionDAL,
		APITokenDAL: apiTokenDAL,
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
