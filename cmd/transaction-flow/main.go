package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	handler "github.com/VieiraVitor/transaction-flow/internal/api/handlers"
	"github.com/VieiraVitor/transaction-flow/internal/application/usecase"
	accountrepository "github.com/VieiraVitor/transaction-flow/internal/infra/repository/account"
	transactionrepository "github.com/VieiraVitor/transaction-flow/internal/infra/repository/transaction"
	"github.com/VieiraVitor/transaction-flow/pkg/logger"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
)

func main() {
	logger.InitLogger()
	logger.Logger.Info("Starting...")
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"db",
		5432,
		"postgres",
		"postgres",
		"transactions",
	)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		logger.Logger.Error("Failed to connect to database", slog.String("error", err.Error()))
		log.Fatal(err)
	}

	defer db.Close()

	accountRepo := accountrepository.NewAccountRepository(db)
	transactionRepo := transactionrepository.NewTransactionRepository(db)

	accountUseCase := usecase.NewAccountUseCase(accountRepo)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo)

	handlers := handler.NewHandlers(
		accountUseCase,
		transactionUseCase,
	)
	routes := handlers.NewRoutes()
	logger.Logger.Info("Server started", slog.String("url", "http://localhost:8080"))
	log.Fatal(http.ListenAndServe(":8080", routes))
}
