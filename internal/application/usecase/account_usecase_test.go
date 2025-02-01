package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAccountUseCase_CreateAccount_WhenValidInput_ShouldReturnID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)
	accountUsecase := NewAccountUseCase(mockRepo)

	ctx := context.Background()
	documentNumber := "123456789"
	expectedID := int64(1)

	mockRepo.EXPECT().
		CreateAccount(gomock.Any(), gomock.Any()).
		Return(expectedID, nil)

	// Act
	id, err := accountUsecase.CreateAccount(ctx, documentNumber)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func TestAccountUseCase_CreateAccount_WhenFailedToCreateAccount_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)
	accountUsecase := NewAccountUseCase(mockRepo)

	ctx := context.Background()
	documentNumber := "123456789"
	expectedError := errors.New("failed to create account")

	mockRepo.EXPECT().
		CreateAccount(gomock.Any(), gomock.Any()).
		Return(int64(0), expectedError)

	// Act
	id, err := accountUsecase.CreateAccount(ctx, documentNumber)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Equal(t, "failed to create account", err.Error())
}

func TestAccountUseCase_GetAccount_WhenValidInput_ShouldReturnAccount(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)
	accountUsecase := NewAccountUseCase(mockRepo)
	ctx := context.Background()
	accountExpected := &domain.Account{
		ID:             int64(1),
		DocumentNumber: "123456789",
		CreatedAt:      time.Now(),
	}

	mockRepo.EXPECT().
		GetAccount(gomock.Any(), accountExpected.ID).
		Return(accountExpected, nil)

	// Act
	account, err := accountUsecase.GetAccount(ctx, accountExpected.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, account.ID, accountExpected.ID)
	assert.Equal(t, account.DocumentNumber, accountExpected.DocumentNumber)
	assert.Equal(t, account.CreatedAt, accountExpected.CreatedAt)
}

func TestAccountUseCase_GetAccount_WhenFailedToGetAccount_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)
	accountUsecase := NewAccountUseCase(mockRepo)

	ctx := context.Background()
	accountID := int64(1)
	expectedError := errors.New("failed to get account")

	mockRepo.EXPECT().
		GetAccount(gomock.Any(), gomock.Any()).
		Return(nil, expectedError)

	// Act
	result, err := accountUsecase.GetAccount(ctx, accountID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}
