package dto

import (
	"errors"
)

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type CreateAccountResponse struct {
	ID int64 `json:"id"`
}

type GetAccountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

func (c *CreateAccountRequest) Validate() error {
	if c.DocumentNumber == "" {
		return errors.New("document_number is mandatory")
	}

	return nil
}
