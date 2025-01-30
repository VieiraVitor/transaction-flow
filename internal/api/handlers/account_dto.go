package handler

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type GetAccountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
