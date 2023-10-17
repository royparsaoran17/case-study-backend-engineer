package customer

import (
	"context"
	"manage-se/internal/common"
	"manage-se/internal/provider"
	"manage-se/internal/provider/auth"

	"github.com/pkg/errors"
)

type service struct {
	provider *provider.Provider
}

func NewService(provider *provider.Provider) Customer {
	return &service{provider: provider}
}

func (s *service) GetAllCustomer(ctx context.Context, meta *common.Metadata) ([]auth.Customer, error) {
	customers, err := s.provider.Auth.GetListCustomers(ctx, meta)
	if err != nil {
		return nil, errors.Wrap(err, "getting all customers ")
	}

	return customers, nil
}
