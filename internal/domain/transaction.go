package domain

import "time"

type Transaction struct {
	id              int64
	accountID       int64
	operationTypeID OperationType
	amount          float64
	eventDate       time.Time
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

func NewTransaction(accountID int64, operationType OperationType, amount float64, eventDate ...time.Time) Transaction {
	eDate := time.Now()
	if len(eventDate) > 0 {
		eDate = eventDate[0]
	}
	return Transaction{
		accountID:       accountID,
		operationTypeID: operationType,
		amount:          amount,
		eventDate:       eDate,
	}
}

func (t *Transaction) ID() int64 {
	return t.id
}

func (t *Transaction) AccountID() int64 {
	return t.accountID
}

func (t *Transaction) OperationTypeID() OperationType {
	return t.operationTypeID
}

func (t *Transaction) Amount() float64 {
	return t.amount
}

func (t *Transaction) EventDate() time.Time {
	return t.eventDate
}
