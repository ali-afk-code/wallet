package wallet

import (
	"errors"

	"github.com/ali-afk-code/wallet/pkg/types"
)

var ErrorAccNotFound = errors.New("account not found")

type Service struct {
	nextAccountId int64
	accounts      []*types.Account
	// payments      []*types.Payment
}

func (s *Service) FindAccountById(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrorAccNotFound
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, errors.New("account already existed")
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
