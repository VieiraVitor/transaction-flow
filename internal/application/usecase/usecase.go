package usecase

import (
	"context"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, documentNumber string) (int64, error)
	GetAccount(ctx context.Context, accountID int64) (*domain.Account, error)
}

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, accountID int64, operationTypeID int, amount float64) (int64, error)
}
