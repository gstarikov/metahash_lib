package metahash_lib

import (
	"math/big"
	"testing"
)

func TestNetwork(t *testing.T) {
	testAddress := Address("0x00072a082d1efe1f2eed19a1f60007fd3b39d1344dc3e6f5f2")

	mk, _ := NewKey()
	mn, _ := NewMetahashNetwork(mk, DevNetwork)
	hash, err := mn.Transaction(&Transaction{
		To:    testAddress,
		Value: big.NewInt(666),
		//Fee:   big.NewInt(13),
		Nonce: big.NewInt(1),
		//Data:  "",
	})
	if err != nil {
		t.Errorf("Transaction err -> %s", err)
	}
	t.Logf("hash[%s]", hash)

	hr, err := mn.GetTx(hash)
	if err != nil || hr == nil {
		t.Errorf("Tx err -> %v", err)
	}
	t.Logf("Tx[%+v]", hr)

	bal, err := mn.Balance(testAddress)
	if err != nil || bal == nil {
		t.Errorf("Balance err -> %v", err)
	}
	t.Logf("balance[%+v]", bal)

	err = mn.Add(testAddress)
	if err != nil {
		t.Errorf("Add error -> %s", err)
	}

	bal, err = mn.Balance(testAddress)
	if err != nil || bal == nil {
		t.Errorf("Balance err -> %v", err)
	}
	t.Logf("balance[%+v]", bal)

	hist, err := mn.History(testAddress)
	if err != nil || hist == nil {
		t.Errorf("History err -> %v", err)
	}
	t.Logf("history[%+v]", hist)

	t.Fail()
}
