package usecase

import (
	"context"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/internal/infra/repository"
)

type accountUseCase struct {
	repo repository.AccountRepository
}

func NewAccountUseCase(repo repository.AccountRepository) AccountUseCase {
	return &accountUseCase{
		repo: repo,
	}
}

func (a *accountUseCase) CreateAccount(ctx context.Context, documentNumber string) (int64, error) {
	account := domain.NewAccount(documentNumber)
	return a.repo.CreateAccount(ctx, account)
}

func (a *accountUseCase) GetAccount(ctx context.Context, accountID int64) (*domain.Account, error) {
	return a.repo.GetAccount(ctx, accountID)
}
