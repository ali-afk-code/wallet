package wallet

import (
	"testing"
)

func TestService_FindAccountById_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+998906657700")
	svc.RegisterAccount("+998906650077")
	acc, err := svc.FindAccountById(2)
	if err != nil {
		t.Errorf("%v", err)
	}
	if acc.ID != 2 {
		t.Error("func working incorrectly")
	}
}
func TestService_FindAccountById_failure(t *testing.T) { //how to test failure
	svc := &Service{}
	svc.RegisterAccount("+998906657700")
	svc.RegisterAccount("+998906650077")
	_, err := svc.FindAccountById(3)
	if err == nil {
		t.Errorf("%v", err)
	}
}
