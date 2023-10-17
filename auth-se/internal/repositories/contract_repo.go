package repositories

import (
	"auth-se/internal/common"
	"auth-se/internal/entity"
	"context"

	"auth-se/internal/presentations"
)

type Role interface {
	CreateRole(ctx context.Context, input presentations.RoleCreate) error
	UpdateRole(ctx context.Context, roleID string, input presentations.RoleUpdate) error
	FindRoleByID(ctx context.Context, roleID string) (*entity.Role, error)
	GetAllRole(ctx context.Context) ([]entity.Role, error)
}

type Customer interface {
	CreateCustomer(ctx context.Context, input presentations.CustomerCreate) error
	FindCustomerByPhonePassword(ctx context.Context, phone, password string) (*entity.CustomerDetail, error)
	FindCustomerByPhone(ctx context.Context, phone string) (*entity.CustomerDetail, error)
	UpdateCustomer(ctx context.Context, roleID string, input presentations.CustomerUpdate) error
	GetAllCustomer(ctx context.Context, meta *common.Metadata) ([]entity.Customer, error)
	FindCustomerByID(ctx context.Context, customerID string) (*entity.CustomerDetail, error)
}
