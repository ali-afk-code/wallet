package wallet

import (
	"errors"

	"github.com/ali-afk-code/wallet/pkg/types"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrAmountMustBePositive = errors.New("ammount must be positive")
var ErrPhoneRegistered = errors.New("phone already registered")

type Service struct {
	nextAccountId int64
	accounts      []*types.Account
	// payments      []*types.Payment
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}
func (s *Service) Pay() {

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
