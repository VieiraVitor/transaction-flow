package dto

import (
	"errors"
)

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" example:"1234567890"`
}

type CreateAccountResponse struct {
	ID int64 `json:"id" example:"1"`
}

type GetAccountResponse struct {
	AccountID      int64  `json:"account_id" example:"1"`
	DocumentNumber string `json:"document_number" example:"1234567890"`
}

func (c *CreateAccountRequest) Validate() error {
	if c.DocumentNumber == "" {
		return errors.New("document_number is mandatory")
	}

	return nil
}
