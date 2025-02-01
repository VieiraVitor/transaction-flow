package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTransactionUseCase_CreateTransaction_WhenValidInput_ShouldReturnID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	transactionUsecase := NewTransactionUseCase(mockRepo)

	ctx := context.Background()
	expectedID := int64(1)

	mockRepo.EXPECT().
		CreateTransaction(gomock.Any(), gomock.Any()).
		Return(expectedID, nil)

	// Act
	id, err := transactionUsecase.CreateTransaction(ctx, int64(2), int(domain.Pagamento), 100)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func TestTransactionUseCase_CreateTransaction_WhenFailedToCreateTransaction_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	transactionUsecase := NewTransactionUseCase(mockRepo)

	ctx := context.Background()
	expectedError := errors.New("failed to create transaction")

	mockRepo.EXPECT().
		CreateTransaction(gomock.Any(), gomock.Any()).
		Return(int64(0), expectedError)

	// Act
	id, err := transactionUsecase.CreateTransaction(ctx, int64(1), int(domain.Saque), -10)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Equal(t, "failed to create transaction", err.Error())
}

func TestTransactionUseCase_CreateTransaction_WhenIsPurchaseOrdWithdraw_ShouldEnsureAmountIsNegative(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	transactionUsecase := NewTransactionUseCase(mockRepo)
	ctx := context.Background()

	expectedAmount := -100.50
	mockRepo.EXPECT().
		CreateTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, transaction domain.Transaction) (int64, error) {
			assert.Equal(t, expectedAmount, transaction.Amount)
			return int64(1), nil
		})

	// Act
	id, err := transactionUsecase.CreateTransaction(ctx, 1, int(domain.CompraAVista), 100.50)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestTransactionUseCase_CreateTransaction_WhenIsPayment_ShouldEnsureAmountIsPositive(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	transactionUsecase := NewTransactionUseCase(mockRepo)
	ctx := context.Background()

	expectedAmount := 100.50
	mockRepo.EXPECT().
		CreateTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, transaction domain.Transaction) (int64, error) {
			assert.Equal(t, expectedAmount, transaction.Amount)
			return int64(1), nil
		})

	// Act
	id, err := transactionUsecase.CreateTransaction(ctx, 1, int(domain.Pagamento), -100.50)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestTransactionUseCase_CreateTransaction_WhenInvalidInput_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	transactionUsecase := NewTransactionUseCase(mockRepo)

	ctx := context.Background()
	expectedError := fmt.Errorf("invalid operation type: %v", 10)

	// Act
	id, err := transactionUsecase.CreateTransaction(ctx, int64(1), 10, 50)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Equal(t, expectedError, err)
}
