package wallet

import (
	"errors"

	"github.com/ali-afk-code/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrAmountMustBePositive = errors.New("ammount must be positive")
var ErrPhoneRegistered = errors.New("phone already registered")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("not found payment with thi id")

type Service struct {
	nextAccountId int64
	accounts      []*types.Account
	payments      []*types.Payment
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountId++
	acc := &types.Account{
		ID:    s.nextAccountId,
		Phone: phone,
	}
	s.accounts = append(s.accounts, acc)
	return acc, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	var acc *types.Account

	for _, account := range s.accounts {
		if accountID == account.ID {
			acc = account
		}
	}
	if acc == nil {
		return ErrAccountNotFound
	}
	acc.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var acc *types.Account

	for _, account := range s.accounts {
		if accountID == account.ID {
			acc = account
			break
		}
	}
	if acc == nil {
		return nil, ErrAccountNotFound
	}
	if acc.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	acc.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	if payment.Status == types.PaymentStatusInProgress {
		payment.Status = types.PaymentStatusFail
		acc, err := s.FindAccountByID(payment.AccountID)
		if err != nil {
			return err
		}
		acc.Balance += payment.Amount
	}
	return nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound

}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	paymentCopy := types.Payment{
		ID:        uuid.NewString(),
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Category:  payment.Category,
		Status:    payment.Status,
	}
	acc, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil, err
	}
	acc.Balance -= payment.Amount
	return &paymentCopy, nil
}
