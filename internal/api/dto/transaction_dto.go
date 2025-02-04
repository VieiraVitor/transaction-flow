package dto

import "errors"

type CreateTransactionRequest struct {
	AccountID       int64   `json:"account_id" example:"1"`
	OperationTypeID int     `json:"operation_type_id" example:"4"`
	Amount          float64 `json:"amount" example:"100"`
}

type CreateTransactionResponse struct {
	ID int64 `json:"id" example:"1"`
}

func (c *CreateTransactionRequest) Validate() error {
	if c.AccountID == 0 {
		return errors.New("accountID is mandatory")
	}

	if c.OperationTypeID == 0 {
		return errors.New("operationTypeID is mandatory")
	}

	if c.Amount == 0 {
		return errors.New("amount is mandatory")
	}

	return nil
}
