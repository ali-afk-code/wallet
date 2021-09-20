package wallet

import (
	"testing"
)

func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+998906657700")
	svc.RegisterAccount("+998906650077")
	acc, err := svc.FindAccountByID(2)
	if err != nil {
		t.Errorf("%v", err)
	}
	if acc.ID != 2 {
		t.Error("func working incorrectly")
	}
}
func TestService_FindAccountByID_failure(t *testing.T) { //how to test failure
	svc := &Service{}
	svc.RegisterAccount("+998906657700")
	svc.RegisterAccount("+998906650077")
	_, err := svc.FindAccountByID(3)
	if err == nil {
		t.Errorf("%v", err)
	}
}
