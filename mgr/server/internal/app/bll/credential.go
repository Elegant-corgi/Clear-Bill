package bll

import (
	"context"
	"errors"
	"strings"
	"time"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/dal"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
)

type CredentialService struct {
	CredentialDAL *dal.CredentialDAL
}

func NewCredentialService(credentialDAL *dal.CredentialDAL) *CredentialService {
	return &CredentialService{CredentialDAL: credentialDAL}
}

func (s *CredentialService) CreateCredential(ctx context.Context, actor *dbmodel.User, req *vo.CreateCredentialReq) (*vo.Credential, error) {
	typ, err := normalizeCredentialType(req.Type)
	if err != nil {
		return nil, err
	}

	item, err := buildCredential(actor.ID, typ, req.Name, nil)
	if err != nil {
		return nil, err
	}

	if err := s.CredentialDAL.Create(ctx, item); err != nil {
		return nil, err
	}

	return toCredentialVO(item, true), nil
}

func (s *CredentialService) RotateCredential(ctx context.Context, actor *dbmodel.User, id uint, req *vo.RotateCredentialReq) (*vo.Credential, error) {
	existing, err := s.CredentialDAL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canAccessCredential(actor, existing) {
		return nil, errors.New("permission denied")
	}

	newName := strings.TrimSpace(req.Name)
	if newName == "" {
		newName = existing.Name
	}

	rotatedAt := time.Now()
	existing.Status = dbmodel.CredentialStatusRotated
	existing.RotatedAt = &rotatedAt
	if err := s.CredentialDAL.Update(ctx, existing); err != nil {
		return nil, err
	}

	item, err := buildCredential(actor.ID, existing.Type, newName, &existing.ID)
	if err != nil {
		existing.Status = dbmodel.CredentialStatusActive
		existing.RotatedAt = nil
		_ = s.CredentialDAL.Update(ctx, existing)
		return nil, err
	}

	if err := s.CredentialDAL.Create(ctx, item); err != nil {
		existing.Status = dbmodel.CredentialStatusActive
		existing.RotatedAt = nil
		_ = s.CredentialDAL.Update(ctx, existing)
		return nil, err
	}

	return toCredentialVO(item, true), nil
}

func (s *CredentialService) ListCredentials(ctx context.Context, actor *dbmodel.User, req *vo.ListCredentialReq) (*vo.PageResult[vo.Credential], error) {
	pageReq := req.PageReq.Normalize()
	typ, err := normalizeCredentialTypeOrEmpty(req.Type)
	if err != nil {
		return nil, err
	}
	status := strings.TrimSpace(req.Status)
	if status != "" {
		switch status {
		case dbmodel.CredentialStatusActive, dbmodel.CredentialStatusRotated, dbmodel.CredentialStatusRevoked:
		default:
			return nil, errors.New("invalid credential status")
		}
	}

	var userID *uint
	if actor.Role != dbmodel.RoleSysadmin {
		userID = &actor.ID
	}

	items, query, total, err := s.CredentialDAL.List(ctx, userID, typ, status, pageReq.Page, pageReq.PageSize)
	if err != nil {
		return nil, err
	}

	result := make([]vo.Credential, 0, len(items))
	for i := range items {
		result = append(result, *toCredentialVO(&items[i], false))
	}

	return &vo.PageResult[vo.Credential]{
		List:     result,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

func (s *CredentialService) GetCredential(ctx context.Context, actor *dbmodel.User, id uint) (*vo.Credential, error) {
	item, err := s.CredentialDAL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canAccessCredential(actor, item) {
		return nil, errors.New("permission denied")
	}
	return toCredentialVO(item, false), nil
}

func (s *CredentialService) DeleteCredential(ctx context.Context, actor *dbmodel.User, id uint) error {
	item, err := s.CredentialDAL.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !canAccessCredential(actor, item) {
		return errors.New("permission denied")
	}
	return s.CredentialDAL.Delete(ctx, id)
}

func buildCredential(userID uint, typ, name string, parentID *uint) (*dbmodel.UserCredential, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("credential name is required")
	}

	now := time.Now()
	credential := &dbmodel.UserCredential{
		UserID:   userID,
		Type:     typ,
		Name:     name,
		Status:   dbmodel.CredentialStatusActive,
		ParentID: parentID,
	}

	switch typ {
	case dbmodel.CredentialTypeAKSK:
		accessKey, err := generateCredentialValue("ak_")
		if err != nil {
			return nil, err
		}
		secretKey, err := generateCredentialValue("sk_")
		if err != nil {
			return nil, err
		}
		credential.AccessKey = &accessKey
		credential.SecretKey = &secretKey
	case dbmodel.CredentialTypeToken:
		token, err := generateToken()
		if err != nil {
			return nil, err
		}
		credential.Token = &token
	default:
		return nil, errors.New("invalid credential type")
	}

	credential.CreatedAt = now
	credential.UpdatedAt = now
	return credential, nil
}

func normalizeCredentialType(raw string) (string, error) {
	typ := strings.ToLower(strings.TrimSpace(raw))
	switch typ {
	case dbmodel.CredentialTypeAKSK, dbmodel.CredentialTypeToken:
		return typ, nil
	default:
		return "", errors.New("invalid credential type")
	}
}

func normalizeCredentialTypeOrEmpty(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}
	return normalizeCredentialType(raw)
}

func canAccessCredential(actor *dbmodel.User, credential *dbmodel.UserCredential) bool {
	if actor == nil || credential == nil {
		return false
	}
	if actor.Role == dbmodel.RoleSysadmin {
		return true
	}
	return actor.ID == credential.UserID
}

func toCredentialVO(item *dbmodel.UserCredential, revealSecrets bool) *vo.Credential {
	result := &vo.Credential{
		ID:         item.ID,
		UserID:     item.UserID,
		Type:       item.Type,
		Name:       item.Name,
		Status:     item.Status,
		ParentID:   item.ParentID,
		ExpiresAt:  item.ExpiresAt,
		RotatedAt:  item.RotatedAt,
		LastUsedAt: item.LastUsedAt,
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
	}
	if item.AccessKey != nil {
		result.AccessKeyPreview = maskCredentialValue(*item.AccessKey)
		if revealSecrets {
			result.AccessKey = *item.AccessKey
		}
	}
	if item.SecretKey != nil && revealSecrets {
		result.SecretKey = *item.SecretKey
	}
	if item.Token != nil {
		result.TokenPreview = maskCredentialValue(*item.Token)
		if revealSecrets {
			result.Token = *item.Token
		}
	}
	return result
}

func maskCredentialValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}
	return value[:4] + "****" + value[len(value)-4:]
}

func generateCredentialValue(prefix string) (string, error) {
	value, err := generateToken()
	if err != nil {
		return "", err
	}
	return prefix + value, nil
}
