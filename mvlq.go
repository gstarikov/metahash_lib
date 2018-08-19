package metahash_lib

import (
	"bytes"
	"encoding/hex"
	"math/big"
)

type MVLQ interface {
	GetData() string
	Append(*big.Int) error
	AppendString(string)
}

func NewMVLQ() MVLQ {
	return new(mvlqImpV1)
}

type mvlqImpV1 struct {
	buffer bytes.Buffer
}

func (t *mvlqImpV1) GetData() string {
	return t.buffer.String()
}

/*
 first byte      value
 0-249 	         the same number
 250 (0xfa) 	 as uint16
 251 (0xfb) 	 as uint32
 252 (0xfc) 	 as uint64
 253 (0xfd) 	 as uint128
 254 (0xfe) 	 as uint256
 255 (0xff) 	 as uint512
*/

func GenerateMVLQ(number *big.Int) (bytes.Buffer, error) {
	var ret bytes.Buffer

	bitLen := number.BitLen()
	var b, p int

	switch {
	case number.Sign() < 0:
		return ret, &ErrorNegativeNumber{}
	case number.Sign() == 0:
		ret.WriteByte(0)
		return ret, nil
	case number.Cmp(big.NewInt(249)) <= 0:
		ret.WriteByte(number.Bytes()[0])
		return ret, nil
	case bitLen > 512:
		return ret, &ErrorTooBigNumber{}
	default:
		for b, p = 16, 250; b < 512 && b < bitLen; b, p = b*2, p+1 {
		}
		//write tag
		ret.WriteByte(byte(p))
		//write number, BE as LE
		numberBytes := number.Bytes()
		for i := (bitLen - 1) / 8; i >= 0; i-- {
			ret.WriteByte(numberBytes[i])
		}
		//write leading zeroes (eg 300 bits number & 512 necessary by format)
		zeroBytes := (b - bitLen) / 8
		for i := 0; i < zeroBytes; i++ {
			ret.WriteByte(0)
		}
		return ret, nil
	}
}

func (t *mvlqImpV1) Append(number *big.Int) error {
	mvlq, err := GenerateMVLQ(number)
	if err != nil {
		return err
	}
	t.buffer.WriteString(hex.EncodeToString(mvlq.Bytes()))
	return nil
}

func (t *mvlqImpV1) AppendString(str string) {
	if str == "" {
		str = "00"
	}
	t.buffer.WriteString(str)
}

type ErrorNegativeNumber struct{}

func (e *ErrorNegativeNumber) Error() string {
	return "ErrorNegativeNumber"
}

type ErrorTooBigNumber struct{}

func (e *ErrorTooBigNumber) Error() string {
	return "ErrorTooBigNumber"
}

func SignTransaction(tr *Transaction, mk MetahashKey) (string, error) {
	var err error
	mlvq := NewMVLQ()
	mlvq.AppendString(tr.To)
	if err := mlvq.Append(tr.Value); err != nil {
		return "", err
	}
	if err := mlvq.Append(tr.Fee); err != nil {
		return "", err
	}
	if err := mlvq.Append(tr.Nonce); err != nil {
		return "", err
	}
	mlvq.AppendString(tr.Data)

	sign, err := mk.Sign([]byte(mlvq.GetData()))
	if err != nil {
		return "", err
	}
	return sign, nil
}
