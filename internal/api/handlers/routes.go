package handler

import (
	"github.com/VieiraVitor/transaction-flow/internal/api/middleware"
	"github.com/VieiraVitor/transaction-flow/internal/application/usecase"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	accountHandler     *AccountHandler
	transactionHandler *TransactionHandler
}

func NewHandlers(
	accountUC *usecase.AccountUseCase,
	transactionUC *usecase.TransactionUseCase,
) *Handlers {
	return &Handlers{
		accountHandler:     NewAccountHandler(accountUC),
		transactionHandler: NewTransactionHandler(transactionUC),
	}
}

func (h *Handlers) NewRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.accountHandler.CreateAccount)
		r.Get("/{id}", h.accountHandler.GetAccount)
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", h.transactionHandler.CreateTransaction)
	})

	return r
}
