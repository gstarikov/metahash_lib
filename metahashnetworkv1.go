package metahash_lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type metaHashRequest struct {
	JsonRPC string
	Method  string
	Params  metahashTransaction
}

type metahashTransactionStrings struct {
	To    string
	Value string
	Fee   string
	Nonce string
	Data  string
}

type metaHashResponce struct {
	Result string
	Params string
	Error  string
}

type metahashTransaction struct {
	Transaction *metahashTransactionStrings
	Pubkey      string
	Sign        string
}

type metahashNetworkImpV1 struct {
	mk MetahashKey
	metahashNetworkPublicImpV1
}

func newMetahashNetworkV1(mk MetahashKey, url string) (MetahashNetwork, error) {
	return &metahashNetworkImpV1{mk: mk,
		metahashNetworkPublicImpV1: metahashNetworkPublicImpV1{
			url: url,
		},
	}, nil
}

type metahashNetworkPublicImpV1 struct {
	url string
	mp  MetahashPublic
}

func newMetahashNetworkPublicV1(mp MetahashPublic, url string) (MetahashNetworkPublic, error) {
	return &metahashNetworkPublicImpV1{mp: mp, url: url}, nil
}

type ErrorNetwork struct{}

func (t *ErrorNetwork) Error() string {
	return "ErrorNetwork code != 200"
}

func (t *metahashNetworkImpV1) Transaction(tr *Transaction) (TxHash, error) {
	//log.Printf("")
	sign, err := SignTransaction(tr, t.mk)
	if err != nil {
		return "", err
	}

	req := metaHashRequest{
		JsonRPC: "2.0",
		Method:  "mhc_send",
		Params: metahashTransaction{
			Transaction: &metahashTransactionStrings{
				To:    tr.To,
				Value: tr.Value.String(),
				Fee:   tr.Fee.String(),
				Nonce: tr.Nonce.String(),
				Data:  tr.Data,
			},
			Pubkey: t.mk.Public(),
			Sign:   sign,
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	//log.Printf("%s",string(reqJson))

	resp, err := http.Post(t.url, "application/json", bytes.NewBuffer(reqJson))
	if resp != nil && resp.Body != nil {
		defer io.Copy(ioutil.Discard, resp.Body)
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	//log.Printf("")

	if resp.StatusCode != 200 {
		return "", &ErrorNetwork{}
	}
	//log.Printf("")

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//log.Printf("%s", string(respBody))

	var respStruct metaHashResponce

	err = json.Unmarshal(respBody, &respStruct)
	if err != nil {
		return "", err
	}

	if respStruct.Error != "" || respStruct.Result != "ok" || respStruct.Params == "" {
		return "", fmt.Errorf("%v", respStruct)
	}

	//wow!!
	//log.Printf("")
	return TxHash(respStruct.Params), nil

}

//type MetahashNetwork interface {
//	MetahashNetworkPublic
//	Transaction(Transaction) error
//}
//
//type MetahashNetworkPublic interface {
//	Balance()
//	History()
//}
