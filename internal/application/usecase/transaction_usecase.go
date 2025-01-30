package usecase

import (
	"context"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
	transactionRepository "github.com/VieiraVitor/transaction-flow/internal/infra/repository/transaction"
)

type TransactionUseCase struct {
	repo transactionRepository.Repository
}

func NewTransactionUseCase(repo transactionRepository.Repository) *TransactionUseCase {
	return &TransactionUseCase{
		repo: repo,
	}
}

func (t *TransactionUseCase) CreateTransaction(ctx context.Context, accountID int, operationTypeID int, amount float64) (int, error) {
	transaction := domain.NewTransaction(accountID, operationTypeID, amount)
	return t.repo.CreateTransaction(ctx, transaction)
}
