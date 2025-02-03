package domain

type Transaction struct {
	ID              string        `json:"id"`
	AccountID       int64         `json:"account_id"`
	OperationTypeID OperationType `json:"operation_type_id"`
	Amount          float64       `json:"amount"`
}

type OperationType int

const (
	CompraAVista    OperationType = 1
	CompraParcelada OperationType = 2
	Saque           OperationType = 3
	Pagamento       OperationType = 4
)

func (o OperationType) IsValid() bool {
	return o == CompraAVista || o == CompraParcelada || o == Saque || o == Pagamento
}

func (o OperationType) IsPayment() bool {
	return o == Pagamento
}

func (o OperationType) IsPurchaseOrWithdraw() bool {
	return o == CompraAVista || o == CompraParcelada || o == Saque
}

func NewTransaction(accountID int64, operationType OperationType, amount float64) Transaction {
	return Transaction{
		AccountID:       accountID,
		OperationTypeID: operationType,
		Amount:          amount,
	}
}
