package api

import (
	"github.com/VieiraVitor/transaction-flow/internal/api/handler"
	"github.com/VieiraVitor/transaction-flow/internal/api/middleware"
	"github.com/VieiraVitor/transaction-flow/internal/application/usecase"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handlers struct {
	accountHandler     *handler.AccountHandler
	transactionHandler *handler.TransactionHandler
}

func NewHandlers(
	accountUseCase usecase.AccountUseCase,
	transactionUseCase usecase.TransactionUseCase,
) *Handlers {
	return &Handlers{
		accountHandler:     handler.NewAccountHandler(accountUseCase),
		transactionHandler: handler.NewTransactionHandler(transactionUseCase),
	}
}

func (h *Handlers) NewRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware, middleware.RecoverMiddleware)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.accountHandler.CreateAccount)
		r.Get("/{id}", h.accountHandler.GetAccount)
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", h.transactionHandler.CreateTransaction)
	})

	return r
}
