package metahash_lib

import (
	"fmt"
	"log"
	"math/big"
)

// https://support.metahash.org/hc/ru/articles/360002712193
// https://support.metahash.org/hc/ru/articles/360003271694
// http://developers.metahash.org

type Address string
type PrivateKey string
type PublicKey string
type Sign string

var Logger = (*log.Logger)(nil) //log.New(os.Stderr, "", log.LstdFlags)

type MetahashKey interface {
	MetahashPublic
	Private() PrivateKey
	Sign(data []byte) (Sign, error)
}

func NewKey() (MetahashKey, error) {
	return newKeyV1()
}

func CreateKey(private PrivateKey) (MetahashKey, error) {
	return createKeyV1(private)
}

type MetahashPublic interface {
	Public() PublicKey
	Veriff(data []byte, sign Sign) (bool, error)
	Address() Address
}

func CreatePublic(public PublicKey) (MetahashPublic, error) {
	return createPublicV1(public)
}

type Transaction struct {
	To    Address
	Value *big.Int
	//Fee   *big.Int // not implemented yet
	Nonce *big.Int
	//Data  string // not implemented yet
}

type Balance struct {
	Address       Address
	Received      *big.Int
	Spent         *big.Int
	CountReceived int
	CountSpent    int
	BlockNumber   int
	CurrentBlock  int
}

type HistoryRecs []HistoryRec

type HistoryRec struct {
	From   Address  `json:"from"`
	To     Address  `json:"to"`
	Value  *big.Int `json:"value"`
	TxHash TxHash   `json:"transaction"`
}

type TxData struct {
	Transaction HistoryRec
}

type TxHash string

type MetahashNetwork interface {
	MetahashNetworkPublic
	MetahashNetworkDev
	Transaction(*Transaction) (TxHash, error)
}

type MetahashNetworkPublic interface {
	Balance(Address) (*Balance, error)
	History(Address) (*HistoryRecs, error)
	GetTx(TxHash) (*HistoryRec, error)
}

type MetahashNetworkDev interface {
	Add(Address) error
}

func NewMetahashNetwork(mk MetahashKey, net NetworkType) (MetahashNetwork, error) {
	return newMetahashNetworkV1(mk, net)
}

func NewMetahashNetworkPublic(mp MetahashPublic, net NetworkType) (MetahashNetworkPublic, error) {
	return newMetahashNetworkPublicV1(mp, net)
}

func logPrintf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Output(2, fmt.Sprintf(format, v...))
	}
}
