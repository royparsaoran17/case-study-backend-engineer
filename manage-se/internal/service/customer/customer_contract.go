package customer

import (
	"context"
	"manage-se/internal/common"
	"manage-se/internal/provider/auth"
)

type Customer interface {
	GetAllCustomer(ctx context.Context, meta *common.Metadata) ([]auth.Customer, error)
}
