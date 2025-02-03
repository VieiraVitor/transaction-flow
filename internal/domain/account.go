package domain

import "time"

type Account struct {
	id             int64
	documentNumber string
	createdAt      time.Time
}

func NewAccount(documentNumber string) *Account {
	return &Account{
		documentNumber: documentNumber,
	}
}

func (a *Account) ID() int64 {
	return a.id
}

func (a *Account) DocumentNumber() string {
	return a.documentNumber
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) SetID(id int64) {
	a.id = id
}

func (a *Account) SetCreatedAt(createdAt time.Time) {
	a.createdAt = createdAt
}
