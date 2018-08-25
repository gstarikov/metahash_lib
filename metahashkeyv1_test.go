package metahash_lib

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	os.Exit(m.Run())
}

func TestMetahashKeyImpV1_GenerateRestore(t *testing.T) {
	mk, _ := NewKey()
	priv := mk.Private()
	mkRestored, err := CreateKey(priv)

	if err != nil {
		t.Fatalf("cant CreateKey -> %s", err)
	}

	if mk.Private() != mkRestored.Private() {
		t.Errorf("restore mismatch")
	}
}

func TestMetahashKeyImpV1_Sign(t *testing.T) {
	testPrivKey := PrivateKey("30770201010420e546b527f59adca85be22aef5ffccabe72c0f374b1bd01dbd91f0d74a773cca4a00a06082a8648ce3d030107a14403420004d08b01f54ed31f085ac27718c37dd12d5f17a8ccfbb26f2a973122356a66f2087eb0d9464cebe701ca640258083fe9f6516290a5f06750772b661113ca60f495")
	testPubKey := PublicKey("3059301306072a8648ce3d020106082a8648ce3d03010703420004d08b01f54ed31f085ac27718c37dd12d5f17a8ccfbb26f2a973122356a66f2087eb0d9464cebe701ca640258083fe9f6516290a5f06750772b661113ca60f495")
	//testAddr := "0x0099f4d2c76be3455f402b5d0538d84040c62669d565b26c33"
	//testSig := "304402204f8104138b52812c2765b39133cd97ccbf3919e6616dec2e0ec6b314af4debf202205214ff552455bea437fc4562095ebc3276e4bd47501e7ac70f0b50bd94060639"
	testSigMsg := "test"

	mk, err := CreateKey(testPrivKey)
	if err != nil {
		t.Fatal(err)
	}
	priv := mk.Private()
	publ := mk.Public()
	sign, err := mk.Sign([]byte(testSigMsg))
	if err != nil {
		t.Fatalf("sign error -> %s", err)
	}
	t.Logf("private  [%s]", priv)
	t.Logf("public   [%s]", publ)
	t.Logf("signature[%s]", sign)

	if priv != testPrivKey {
		t.Errorf("private key mismatch.\n has[%s]\nwant[%s]", priv, testPrivKey)
	}

	if publ != testPubKey {
		t.Errorf("public key mismatch\n has[%s]\nwant[%s]", publ, testPubKey)
	}

	veriff, err := mk.Veriff([]byte(testSigMsg), sign)
	if err != nil || !veriff {
		t.Errorf("cant veriff. veriff[%t], err -> %v", veriff, err)
	}

	//t.Fail()
}
