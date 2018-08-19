package metahash_lib

import (
	"math/big"
	"testing"
)

//dev network
// proxy.net-dev.metahash.org

func TestNetwork(t *testing.T) {
	mk, _ := NewKey()
	mn, _ := NewMetahashNetwork(mk, "https://proxy.net-dev.metahash.org")
	hash, err := mn.Transaction(&Transaction{
		To:    "0x0057d0697c8e59859608aaa1c4ce11e9685d3d30b02876a632",
		Value: big.NewInt(666),
		Fee:   big.NewInt(13),
		Nonce: big.NewInt(0),
		Data:  "no data",
	})
	if err != nil {
		t.Fatalf("err -> %s", err)
	}
	t.Logf("hash[%s]", hash)
	t.Fail()
}
