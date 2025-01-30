package repository

import (
	"context"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
)

type Repository interface {
	CreateTransaction(ctx context.Context, transaction domain.Transaction) (int, error)
}
