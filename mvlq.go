package metahash_lib

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"strings"
)

type MVLQ interface {
	GetData() []byte
	Append(*big.Int) error
	AppendBytes([]byte)
	AppendString(string) error
}

func NewMVLQ() MVLQ {
	return new(mvlqImpV1)
}

type mvlqImpV1 struct {
	buffer bytes.Buffer
}

func (t *mvlqImpV1) GetData() []byte {
	return t.buffer.Bytes()
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

func GenerateMVLQ(number *big.Int, buffer *bytes.Buffer) (*bytes.Buffer, error) {
	var ret *bytes.Buffer
	if buffer == nil {
		ret = &bytes.Buffer{}
	} else {
		ret = buffer
	}

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
	_, err := GenerateMVLQ(number, &t.buffer)
	return err
}

func (t *mvlqImpV1) AppendString(str string) error {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return err
	}

	t.AppendBytes(bytes)
	return nil
}

func (t *mvlqImpV1) AppendBytes(bytes []byte) {
	if len(bytes) == 0 {
		t.buffer.WriteByte(0)
	} else {
		t.buffer.Write(bytes)
	}
}

type ErrorNegativeNumber struct{}

func (e *ErrorNegativeNumber) Error() string {
	return "ErrorNegativeNumber"
}

type ErrorTooBigNumber struct{}

func (e *ErrorTooBigNumber) Error() string {
	return "ErrorTooBigNumber"
}

func SignTransaction(tr *Transaction, mk MetahashKey) (Sign, error) {
	var err error
	mlvq := NewMVLQ()
	to, _ := hex.DecodeString(strings.TrimPrefix(string(tr.To), "0x"))
	//to := string(tr.To)
	mlvq.AppendBytes(to)
	if err := mlvq.Append(tr.Value); err != nil {
		return "", err
	}
	if err := mlvq.Append(big.NewInt(0) /*tr.Fee*/); err != nil {
		return "", err
	}
	if err := mlvq.Append(tr.Nonce); err != nil {
		return "", err
	}
	mlvq.AppendBytes(nil /*tr.Data*/)

	mlvqData := mlvq.GetData()

	sign, err := mk.Sign(mlvqData)
	if err != nil {
		return "", err
	}

	return sign, nil
}
