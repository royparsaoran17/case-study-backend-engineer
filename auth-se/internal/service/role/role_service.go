package role

import (
	"context"

	"auth-se/internal/entity"
	"auth-se/internal/presentations"
	"auth-se/internal/repositories"
	"github.com/pkg/errors"
)

type service struct {
	repo *repositories.Repository
}

func NewService(repo *repositories.Repository) Role {
	return &service{repo: repo}
}

func (s *service) GetAllRole(ctx context.Context) ([]entity.Role, error) {
	roles, err := s.repo.Role.GetAllRole(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting all roles on ")
	}

	return roles, nil
}

func (s *service) GetRoleByID(ctx context.Context, roleID string) (*entity.Role, error) {
	roles, err := s.repo.Role.FindRoleByID(ctx, roleID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting role id %s", roleID)
	}

	return roles, nil
}

func (s *service) UpdateRoleByID(ctx context.Context, roleID string, input presentations.RoleUpdate) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation(s) error")
	}

	_, err := s.repo.Role.FindRoleByID(ctx, roleID)
	if err != nil {
		return errors.Wrapf(err, "getting role id %s", roleID)
	}

	if err := s.repo.Role.UpdateRole(ctx, roleID, input); err != nil {
		return errors.Wrap(err, "updating role")

	}

	return nil
}

func (s *service) CreateRole(ctx context.Context, input presentations.RoleCreate) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation(s) error")
	}

	err := s.repo.Role.CreateRole(ctx, input)
	if err != nil {
		return errors.Wrap(err, "creating role")

	}

	return nil
}
