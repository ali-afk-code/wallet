package types

import "github.com/google/uuid"

type Money int64
type Phone string
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

type Payment struct {
	PaymentID uuid.UUID
}
