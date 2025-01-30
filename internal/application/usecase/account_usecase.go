package usecase

import (
	"context"
	"time"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
	accountRepository "github.com/VieiraVitor/transaction-flow/internal/infra/repository/account"
)

type AccountUseCase struct {
	repo accountRepository.Repository
}

func NewAccountUseCase(repo accountRepository.Repository) *AccountUseCase {
	return &AccountUseCase{
		repo: repo,
	}
}

func (a *AccountUseCase) CreateAccount(ctx context.Context, documentNumber string) (int, error) {
	account := domain.NewAccount(documentNumber, time.Now())
	return a.repo.CreateAccount(ctx, account)
}

func (a *AccountUseCase) GetAccount(ctx context.Context, accountID int64) (*domain.Account, error) {
	return a.repo.GetAccount(ctx, accountID)
}
