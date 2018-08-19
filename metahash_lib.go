package metahash_lib

import "math/big"

// https://support.metahash.org/hc/ru/articles/360002712193
// https://support.metahash.org/hc/ru/articles/360003271694
// http://developers.metahash.org

type MetahashKey interface {
	MetahashPublic
	Private() string
	//Public() string
	Sign(data []byte) (string, error)
	//Veriff(data []byte, sign string) (bool,error)
	//Address() string
}

func NewKey() (MetahashKey, error) {
	return newKeyV1()
}

func CreateKey(private string) (MetahashKey, error) {
	return createKeyV1(private)
}

type MetahashPublic interface {
	Public() string
	Veriff(data []byte, sign string) (bool, error)
	Address() string
}

func CreatePublic(public string) (MetahashPublic, error) {
	return createPublicV1(public)
}

type Transaction struct {
	To    string
	Value *big.Int
	Fee   *big.Int
	Nonce *big.Int
	Data  string
}

type TxHash string

type MetahashNetwork interface {
	MetahashNetworkPublic
	Transaction(*Transaction) (TxHash, error)
}

type MetahashNetworkPublic interface {
	//Balance()
	//History()
}

func NewMetahashNetwork(mk MetahashKey, connectionString string) (MetahashNetwork, error) {
	return newMetahashNetworkV1(mk, connectionString)
}

func NewMetahashNetworkPublic(mp MetahashPublic, connectionString string) (MetahashNetworkPublic, error) {
	return newMetahashNetworkPublicV1(mp, connectionString)
}
