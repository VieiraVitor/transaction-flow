package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/internal/infra/repository"
)

type transactionUseCase struct {
	repo repository.TransactionRepository
}

func NewTransactionUseCase(repo repository.TransactionRepository) TransactionUseCase {
	return &transactionUseCase{
		repo: repo,
	}
}

func (t *transactionUseCase) CreateTransaction(ctx context.Context, accountID int64, operationTypeID int, amount float64) (int64, error) {
	operationType := domain.OperationType(operationTypeID)

	if !operationType.IsValid() {
		return 0, fmt.Errorf("invalid operation type: %v", operationType)
	}

	if operationType.IsPayment() && amount < 0 {
		amount = -amount
	}

	if operationType.IsPurchaseOrWithdraw() && amount > 0 {
		amount = -amount
	}

	transaction := domain.NewTransaction(accountID, operationType, amount, time.Now())
	return t.repo.CreateTransaction(ctx, transaction)
}
