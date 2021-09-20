package types

type Money int64
type Phone string
type PaymentCategory string
type PaymentStatus string
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

//
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)
