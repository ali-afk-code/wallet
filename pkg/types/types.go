package types

type Money int64
type Phone string
type Account struct {
	ID    int64
	Phone Phone
	Money Money
}

type Payment struct {
}
