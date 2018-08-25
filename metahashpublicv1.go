package metahash_lib

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"reflect"
)

type metahashPublicImpV1 struct {
	pub *ecdsa.PublicKey
}

func createPublicV1(public PublicKey) (MetahashPublic, error) {
	b, err := hex.DecodeString(string(public))
	if err != nil {
		return nil, err
	}

	pub, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}

	p, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cant cast [%s] to [*ecdsa.PublicKey]", reflect.TypeOf(p))
	}

	return &metahashPublicImpV1{
		pub: p,
	}, nil
}

func (t *metahashPublicImpV1) Public() PublicKey {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(t.pub)
	return PublicKey(hex.EncodeToString(x509EncodedPub))
}

func (t *metahashPublicImpV1) Address() Address {
	return ""
}

func (t *metahashPublicImpV1) Veriff(data []byte, sign Sign) (bool, error) {
	digest := sha256.Sum256(data)

	decoded, err := hex.DecodeString(string(sign))
	if err != nil {
		return false, err
	}

	var signEcdsa ecdsaSignature

	_, err = asn1.Unmarshal(decoded, &signEcdsa)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(t.pub, digest[:], signEcdsa.R, signEcdsa.S), nil
}
