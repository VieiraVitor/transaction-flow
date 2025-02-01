package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction domain.Transaction) (int64, error) {
	query := "INSERT INTO transactions (account_id, operation_type_id, amount, event_date) VALUES($1, $2, $3, NOW()) RETURNING id"

	var id int64
	row := r.db.QueryRow(query, transaction.AccountID, transaction.OperationTypeID, transaction.Amount)
	err := row.Scan((&id))
	if err != nil {
		logger.Logger.Error("Error creating transaction", slog.String("error", err.Error()))
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}
	return id, nil
}
