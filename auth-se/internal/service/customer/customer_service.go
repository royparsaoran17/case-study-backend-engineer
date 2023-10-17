package customer

import (
	"auth-se/internal/common"
	"context"
	"github.com/google/uuid"

	"auth-se/internal/entity"
	"auth-se/internal/presentations"
	"auth-se/internal/repositories"
	"github.com/pkg/errors"
)

type service struct {
	repo *repositories.Repository
}

func NewService(repo *repositories.Repository) Customer {
	return &service{repo: repo}
}

func (s *service) GetAllCustomer(ctx context.Context, meta *common.Metadata) ([]entity.Customer, error) {
	customers, err := s.repo.Customer.GetAllCustomer(ctx, meta)
	if err != nil {
		return nil, errors.Wrap(err, "getting all customers on ")
	}

	return customers, nil
}

func (s *service) GetCustomerByID(ctx context.Context, customerID string) (*entity.CustomerDetail, error) {
	customers, err := s.repo.Customer.FindCustomerByID(ctx, customerID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting customer id %s", customerID)
	}

	return customers, nil
}

func (s *service) UpdateCustomerByID(ctx context.Context, customerID string, input presentations.CustomerUpdate) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation(s) error")
	}

	_, err := s.repo.Customer.FindCustomerByID(ctx, customerID)
	if err != nil {
		return errors.Wrapf(err, "getting customer id %s", customerID)
	}

	_, err = s.repo.Role.FindRoleByID(ctx, input.RoleID)
	if err != nil {
		return errors.Wrap(err, "creating customer")

	}

	if err := s.repo.Customer.UpdateCustomer(ctx, customerID, input); err != nil {
		return errors.Wrap(err, "updating customer")

	}

	return nil
}

func (s *service) CreateCustomer(ctx context.Context, input presentations.CustomerCreate) (*entity.CustomerDetail, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "validation(s) error")
	}

	_, err := s.repo.Role.FindRoleByID(ctx, input.RoleID)
	if err != nil {
		return nil, errors.Wrap(err, "creating customer")

	}

	customerID := uuid.NewString()
	input.ID = customerID
	err = s.repo.Customer.CreateCustomer(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "creating customer")

	}

	customer, err := s.repo.Customer.FindCustomerByID(ctx, customerID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting customer id %s", customerID)
	}

	return customer, nil
}
