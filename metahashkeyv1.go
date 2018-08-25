package metahash_lib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"math/big"
)

type metahashKeyImpV1 struct {
	priv *ecdsa.PrivateKey
}

// https://support.metahash.org/hc/ru/articles/360002712193
func newKeyV1() (MetahashKey, error) {
	curve := elliptic.P256() // secp256r1 by default
	rnd := rand.Reader
	priv, err := ecdsa.GenerateKey(curve, rnd)
	if err != nil {
		return nil, err
	}
	return &metahashKeyImpV1{
		priv: priv,
	}, nil

}

func createKeyV1(private PrivateKey) (MetahashKey, error) {
	decoded, err := hex.DecodeString(string(private))
	if err != nil {
		return nil, err
	}
	key, err := x509.ParseECPrivateKey(decoded)
	if err != nil {
		return nil, err
	}
	return &metahashKeyImpV1{
		priv: key,
	}, nil
}

func (t *metahashKeyImpV1) Private() PrivateKey {
	x509EncodedPriv, _ := x509.MarshalECPrivateKey(t.priv)
	return PrivateKey(hex.EncodeToString(x509EncodedPriv))
}
func (t *metahashKeyImpV1) Public() PublicKey {
	return (&metahashPublicImpV1{pub: &t.priv.PublicKey}).Public()
}

type ecdsaSignature struct {
	R, S *big.Int
}

func (t *metahashKeyImpV1) Sign(data []byte) (Sign, error) {
	digest := sha256.Sum256(data)

	r, s, err := ecdsa.Sign(rand.Reader, t.priv, digest[:])
	if err != nil {
		return "", err
	}

	b, e := asn1.Marshal(ecdsaSignature{r, s})

	return Sign(hex.EncodeToString(b)), e
}

func (t *metahashKeyImpV1) Veriff(data []byte, sign Sign) (bool, error) {
	return (&metahashPublicImpV1{pub: &t.priv.PublicKey}).Veriff(data, sign)
}

func (t *metahashKeyImpV1) Address() Address {
	return (&metahashPublicImpV1{pub: &t.priv.PublicKey}).Address()
}
