package auth

import (
	config "auth-se/internal/appctx"
	"auth-se/internal/common"
	"auth-se/internal/consts"
	"auth-se/internal/entity"
	"auth-se/internal/presentations"
	"auth-se/internal/repositories"
	"auth-se/pkg/jwt"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
)

type service struct {
	repo *repositories.Repository
	cfg  config.Common
}

func NewService(repo *repositories.Repository, cfg config.Common) Auth {
	return &service{repo: repo, cfg: cfg}
}

func (s *service) Verify(input presentations.Verify) (*entity.CustomerDetail, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "validation(s) error")
	}

	claim, err := jwt.VerifyToken(input.Token, s.cfg.JWTKey)
	if err != nil {
		return nil, errors.Wrap(err, "verify token")
	}

	claims := *claim

	var customer entity.CustomerDetail
	jsonb, err := json.Marshal(claims["data"].(map[string]interface{}))
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	if err := json.Unmarshal(jsonb, &customer); err != nil {
		return nil, errors.Wrap(err, "unmarshal")

	}

	return &customer, nil
}

func (s *service) Login(ctx context.Context, input presentations.Login) (*entity.CustomerDetailToken, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "validation(s) error")
	}

	_, err := s.repo.Customer.FindCustomerByPhone(ctx, input.Phone)
	if err != nil {
		return nil, consts.ErrCustomerNotFound
	}

	hashString := common.Hasher(input.Password)
	customers, err := s.repo.Customer.FindCustomerByPhonePassword(ctx, input.Phone, hashString)
	if err != nil {
		return nil, consts.ErrWrongPassword
	}

	token, err := jwt.CreateToken(customers, jwt.RequestJwt{
		ID:        customers.ID.String(),
		JWTKey:    s.cfg.JWTKey,
		CreatedAt: customers.CreatedAt,
	})
	if err != nil {
		return nil, errors.Wrap(err, "generate token")
	}

	return &entity.CustomerDetailToken{
		CustomerDetail: *customers,
		Token:          *token,
	}, nil
}
