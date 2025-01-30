package repository

import (
	"context"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
)

type Repository interface {
	CreateAccount(ctx context.Context, account domain.Account) (int, error)
	GetAccount(ctx context.Context, accountID int64) (*domain.Account, error)
}
