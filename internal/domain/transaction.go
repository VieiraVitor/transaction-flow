package domain

type Transaction struct {
	ID              string  `json:"id"`
	AccountID       int     `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

func NewTransaction(accountID int, operationTypeID int, amount float64) Transaction {
	return Transaction{
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
	}
}
