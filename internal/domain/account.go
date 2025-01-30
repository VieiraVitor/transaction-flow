package domain

import "time"

type Account struct {
	ID             int64     `json:"id"`
	DocumentNumber string    `json:"document_number"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewAccount(documentNumber string, createdAt time.Time) Account {
	return Account{
		DocumentNumber: documentNumber,
		CreatedAt:      createdAt,
	}
}
