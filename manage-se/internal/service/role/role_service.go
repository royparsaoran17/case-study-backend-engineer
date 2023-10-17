package role

import (
	"context"
	"manage-se/internal/provider"
	"manage-se/internal/provider/auth"

	"github.com/pkg/errors"
)

type service struct {
	provider *provider.Provider
}

func NewService(provider *provider.Provider) Role {
	return &service{provider: provider}
}

func (s *service) GetAllRole(ctx context.Context) ([]auth.Role, error) {
	roles, err := s.provider.Auth.GetListRoles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting all roles on ")
	}

	return roles, nil
}
