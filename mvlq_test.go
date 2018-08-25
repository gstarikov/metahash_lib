package metahash_lib

import (
	"encoding/hex"
	"math/big"
	. "testing"
)

func TestNewMVLQ(t *T) {
	cases := []struct {
		to, value, fee, nonce, data, result string
	}{
		{
			"009806da73b1589f38630649bdee48467946d118059efd6aab",
			"126894",
			"55647",
			"255",
			"",
			"009806da73b1589f38630649bdee48467946d118059efd6aabfbaeef0100fa5fd9faff0000",
		},
		{
			"009806da73b1589f38630649bdee48467946d118059efd6aab",
			"0",
			"0",
			"0",
			"",
			"009806da73b1589f38630649bdee48467946d118059efd6aab00000000",
		},
		{
			"009806da73b1589f38630649bdee48467946d118059efd6aab",
			"4294967295",
			"65535",
			"249",
			"",
			"009806da73b1589f38630649bdee48467946d118059efd6aabfbfffffffffafffff900",
		},
		{
			"009806da73b1589f38630649bdee48467946d118059efd6aab",
			"4294967296",
			"65536",
			"250",
			"",
			"009806da73b1589f38630649bdee48467946d118059efd6aabfc0000000001000000fb00000100fafa0000",
		},
	}

	for _, c := range cases {
		mvlq := NewMVLQ()
		mvlq.AppendString(c.to)
		helperAppend(mvlq, c.value)
		helperAppend(mvlq, c.fee)
		helperAppend(mvlq, c.nonce)
		mvlq.AppendString(c.data)

		if z := hex.EncodeToString(mvlq.GetData()); z != c.result {
			t.Errorf("%+v result != %s", c, z)
		}
	}
}

func helperAppend(mvlq MVLQ, val string) {
	var bInt big.Int
	bInt.SetString(val, 10)
	mvlq.Append(&bInt)
}

func TestGenerateMVLQ(t *T) {
	cases := []struct {
		value, result string
	}{
		{"126894", "fbaeef0100"},
		{"55647", "fa5fd9"},
		{"255", "faff00"},
		{"0", "00"},
		{"4294967295", "fbffffffff"},
		{"65535", "faffff"},
		{"249", "f9"},
		{"4294967296", "fc0000000001000000"},
		{"65536", "fb00000100"},
		{"250", "fafa00"},
	}

	for _, c := range cases {
		var bInt big.Int
		bInt.SetString(c.value, 10)
		b, e := GenerateMVLQ(&bInt, nil)
		h := hex.EncodeToString(b.Bytes())
		if e != nil || h != c.result {
			t.Errorf("err -> %v, value[%s] got[%s] want[%s]", e, c.value, h, c.result)
		}
	}
}
