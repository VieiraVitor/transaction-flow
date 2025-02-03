package repository

import (
	"context"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *domain.Account) (int64, error)
	GetAccount(ctx context.Context, accountID int64) (*domain.Account, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction domain.Transaction) (int64, error)
}
