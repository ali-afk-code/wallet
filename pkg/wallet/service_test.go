package wallet

import (
	"reflect"
	"testing"

	"github.com/ali-afk-code/wallet/pkg/types"
	"github.com/google/uuid"
)

var defaultTestAccount = testAccount{
	phone:   "+992000001",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

type testService struct {
	*Service
}
type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, err
	}
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, err
	}
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, err
		}

	}
	return account, payments, nil
}
func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	acc, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, err
	}
	err = s.Deposit(acc.ID, balance)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

// func TestService_FindPaymentByID_success(t *testing.T) {
// 	s := newTestService()
// 	// svc.RegisterAccount("+998906657700")
// 	// svc.RegisterAccount("+998906650077")
// 	acc, err := s.addAccountWithBalance("+9920000001", 10_000_00)
// 	if err != nil {
// 		t.Errorf("%v", err)
// 		return
// 	}
// 	payment, err := s.Pay(acc.ID, 1000_00, "auto")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	got, err := s.FindPaymentByID(payment.ID)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if !reflect.DeepEqual(payment, got) {
// 		t.Errorf("FindPaymentByID(): wrong payment returned =%v", err)
// 		return
// 	}
// }
func TestService_FindPaymenttByID_failure(t *testing.T) { //how to test failure
	s := newTestService()
	acc, err := s.addAccountWithBalance("+9920000001", 10_000_00)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	_, err = s.Pay(acc.ID, 1000_00, "auto")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = s.FindPaymentByID(uuid.NewString())
	if err == nil {
		t.Error(err)
		return
	}
	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound returned = %v", err)
		return
	}
}
func TestService_FindAccountByID_success(t *testing.T) {
	s := newTestService()
	account, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	got, err := s.FindAccountByID(account.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(got, account) {
		t.Errorf("FindAccountByID is working not properly, error: %v", err)
		return
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(payment, got) {
		t.Error(err)
		return
	}
}
func TestService_Reject_success(t *testing.T) {
	s := newTestService()
	acc, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payments[0].Status = types.PaymentStatusInProgress
	balanceAfterReject := acc.Balance + payments[0].Amount
	err = s.Reject(payments[0].ID)
	if err != nil {
		t.Error(err)
		return
	}
	if balanceAfterReject != acc.Balance {
		t.Error("Error in Reject, got ", acc.Balance, "wanted ", balanceAfterReject)
	}
}
func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	acc, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := payments[0]
	balanceAfterReject := acc.Balance - payment.Amount
	got, err := s.Repeat(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if payment.ID == got.ID && payment.Amount != got.Amount && payment.Category != got.Category && balanceAfterReject != acc.Balance {
		t.Errorf("Error in Repeat,(all fields except for ID should be equal) got: %v, %v want: %v, %v", payment, balanceAfterReject, got, acc.Balance)
	}
}

func TestService_Reject_fail_user(t *testing.T) {
	var svc Service
	svc.RegisterAccount("+992000000001")
	account, err := svc.FindAccountByID(1)

	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	payment, err := svc.Pay(account.ID, 10_00, "Cafe")

	if err != nil {
		t.Errorf("method Pay returned not nil error, account => %v", account)
	}

	pay, err := svc.FindPaymentByID(payment.ID)

	if err != nil {
		t.Errorf("method FindPaymentByID returned not nil error, payment => %v", payment)
	}

	err = svc.Reject(pay.ID + "uu")

	if err == nil {
		t.Errorf("method Reject returned not nil error, pay => %v", pay)
	}

}

func TestService_Repeat_success_user(t *testing.T) {
	var svc Service

	account, err := svc.RegisterAccount("+992000000001")

	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	payment, err := svc.Pay(account.ID, 10_00, "Cafe")

	if err != nil {
		t.Errorf("method Pay returned not nil error, account => %v", account)
	}

	pay, err := svc.FindPaymentByID(payment.ID)

	if err != nil {
		t.Errorf("method FindPaymentByID returned not nil error, payment => %v", payment)
	}

	paymentNew, err := svc.Repeat(pay.ID)

	if err != nil {
		t.Errorf("method Repat returned not nil error, paymentNew => %v", paymentNew)
	}

}
func TestService_Favorite_success_user(t *testing.T) {
	var svc Service

	account, err := svc.RegisterAccount("+992000000001")

	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	payment, err := svc.Pay(account.ID, 10_00, "Cafe")

	if err != nil {
		t.Errorf("method Pay returned not nil error, account => %v", account)
	}

	favorite, err := svc.FavoritePayment(payment.ID, "My Favorite")

	if err != nil {
		t.Errorf("method FavoritePayment returned not nil error, favorite => %v", favorite)
	}

	paymentFavorite, err := svc.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("method PayFromFavorite returned not nil error, paymentFavorite => %v", paymentFavorite)
	}

}
