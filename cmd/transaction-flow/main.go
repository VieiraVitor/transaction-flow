package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VieiraVitor/transaction-flow/config"
	_ "github.com/VieiraVitor/transaction-flow/docs"
	"github.com/VieiraVitor/transaction-flow/internal/api/handler"
	"github.com/VieiraVitor/transaction-flow/internal/application/usecase"
	"github.com/VieiraVitor/transaction-flow/internal/infra/database"
	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
	"github.com/VieiraVitor/transaction-flow/internal/infra/repository"
	_ "github.com/lib/pq"
)

// @title Transaction Flow API
// @version 1.0
// @description API for account and transaction management
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	logger.InitLogger()

	cfg := config.LoadConfig()

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	accountUseCase := usecase.NewAccountUseCase(accountRepo)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo)

	handlers := handler.NewHandlers(
		accountUseCase,
		transactionUseCase,
	)
	routes := handlers.NewRoutes()

	server := &http.Server{
		Addr:    cfg.AppPort,
		Handler: routes,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Logger.Info("Server started", "url", "http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Error("Failed to start server", "error", err.Error())
		}
	}()

	<-stop
	logger.Logger.Info("Signal received. Stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Logger.Error("Failed to stop server", "error", err.Error())
	} else {
		logger.Logger.Info("Sever finished successfully")
	}
}
