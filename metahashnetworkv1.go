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
	JsonRPC string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  metahashTransaction `json:"params"`
}

type metahashTransactionStrings struct {
	To    Address `json:"to"`
	Value string  `json:"value"`
	Fee   string  `json:"fee"`
	Nonce string  `json:"nonce"`
	Data  string  `json:"data"`
}

type metaHashResponceTransaction struct {
	Result string
	Params string
	Error  string
}

type metahashTransaction struct {
	metahashTransactionStrings
	Pubkey string `json:"pubkey"`
	Sign   string `json:"sign"`
}

type metahashNetworkImpV1 struct {
	mk MetahashKey
	metahashNetworkPublicImpV1
}

func newMetahashNetworkV1(mk MetahashKey, net NetworkType) (MetahashNetwork, error) {
	return &metahashNetworkImpV1{mk: mk,
		metahashNetworkPublicImpV1: metahashNetworkPublicImpV1{
			net: net,
		},
	}, nil
}

type metahashNetworkPublicImpV1 struct {
	net NetworkType
	mp  MetahashPublic
}

func newMetahashNetworkPublicV1(mp MetahashPublic, net NetworkType) (MetahashNetworkPublic, error) {
	return &metahashNetworkPublicImpV1{mp: mp, net: net}, nil
}

type ErrorNetwork struct{}

func (t *ErrorNetwork) Error() string {
	return "ErrorNetwork code != 200"
}

type ErrorNetworkUnreachable struct{}

func (t *ErrorNetworkUnreachable) Error() string {
	return "ErrorNetworkUnreachable"
}

type ErrorNetworkUnsupportedMethod struct{}

func (t *ErrorNetworkUnsupportedMethod) Error() string {
	return "ErrorNetworkUnsupportedMethod"
}

type sendMethod int

const (
	post sendMethod = iota
	get             = iota
)

func (t sendMethod) String() string {
	switch t {
	case post:
		return "POST"
	case get:
		return "GET"
	}
	return "unknown method"
}

func helperSend(urls []string, req []byte, method sendMethod) ([]byte, error) {
	for _, url := range urls {
		var resp *http.Response
		var err error
		//logPrintf("trying [%s]",url)
		switch method {
		case post:
			resp, err = http.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(req))
		case get:
			resp, err = http.Get(url)
		default:
			panic("unsupported method")
		}
		if resp != nil && resp.Body != nil {
			defer io.Copy(ioutil.Discard, resp.Body)
			defer resp.Body.Close()
		}
		if err != nil {
			logPrintf("url[%s] err[%v]", url, err)
			continue
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			logPrintf("url[%s] Code[%d] Response[%s] Request[%s]", url, resp.StatusCode, respBody, req)
			return nil, &ErrorNetwork{}
		}

		logPrintf("%s url[%s], request -> [%s] response -> [%s]", method, url, string(req), string(respBody))
		return respBody, nil
	}
	return nil, &ErrorNetworkUnreachable{}
}

func (t *metahashNetworkImpV1) Transaction(tr *Transaction) (TxHash, error) {
	sign, err := SignTransaction(tr, t.mk)
	if err != nil {
		return "", err
	}

	req := metaHashRequest{
		JsonRPC: "2.0",
		Method:  "mhc_send",
		Params: metahashTransaction{
			metahashTransactionStrings: metahashTransactionStrings{
				To:    tr.To,
				Value: tr.Value.String(),
				//Fee:   tr.Fee.String(),
				Nonce: tr.Nonce.String(),
				//Data:  tr.Data,
			},
			Pubkey: string(t.mk.Public()),
			Sign:   string(sign),
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	url, _ := t.net.ProxyUrl("")

	respBody, err := helperSend(url, reqJson, post)
	if err != nil || respBody == nil {
		return "", err
	}

	var respStruct metaHashResponceTransaction

	err = json.Unmarshal(respBody, &respStruct)
	if err != nil {
		return "", err
	}

	if respStruct.Error != "" || respStruct.Result != "ok" || respStruct.Params == "" {
		return "", fmt.Errorf("%v", respStruct.Error)
	}

	//wow!!
	return TxHash(respStruct.Params), nil
}

type metahashRequestHeader struct {
	Id int `json:"id"`
}

type metahashResponceHeader metahashRequestHeader

type metahashRequestAddress struct {
	Address Address `json:"address"`
}

type metahashRequestTxHash struct {
	Hash TxHash `json:"hash"`
}

type metahashRequestBalance struct {
	metahashRequestHeader
	Params metahashRequestAddress `json:"params"`
}

type metahashResponceBalance struct {
	metahashResponceHeader
	Result *Balance `json:"result"`
}

func (t *metahashNetworkPublicImpV1) Balance(addr Address) (*Balance, error) {
	req := metahashRequestBalance{
		metahashRequestHeader: metahashRequestHeader{
			Id: 1,
		},
		Params: metahashRequestAddress{
			Address: addr,
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url, _ := t.net.TorrentUrl("fetch-balance")

	resp, err := helperSend(url, reqJson, post)
	if err != nil || resp == nil {
		return nil, err
	}

	var respStruct metahashResponceBalance

	err = json.Unmarshal(resp, &respStruct)
	if err != nil {
		return nil, err
	}

	return respStruct.Result, nil
}

type metahashRequestHistory struct {
	metahashRequestHeader
	Params metahashRequestAddress `json:"params"`
}

type metahashResponceHistory struct {
	metahashResponceHeader
	Result HistoryRecs `json:"result"`
}

func (t *metahashNetworkPublicImpV1) History(addr Address) (*HistoryRecs, error) {
	req := metahashRequestBalance{
		metahashRequestHeader: metahashRequestHeader{
			Id: 1,
		},
		Params: metahashRequestAddress{
			Address: addr,
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url, _ := t.net.TorrentUrl("fetch-history")

	resp, err := helperSend(url, reqJson, post)
	if err != nil || resp == nil {
		return nil, err
	}

	var respStruct metahashResponceHistory

	err = json.Unmarshal(resp, &respStruct)
	if err != nil {
		return nil, err
	}

	return &respStruct.Result, nil
}

type metahashRequestGetTx struct {
	metahashRequestHeader
	Params metahashRequestTxHash `json:"params"`
}

type metahashResponceGetTx struct {
	metahashResponceHeader
	Result *HistoryRec `json:"result"`
}

func (t *metahashNetworkPublicImpV1) GetTx(tx TxHash) (*HistoryRec, error) {
	req := metahashRequestGetTx{
		metahashRequestHeader: metahashRequestHeader{
			Id: 1,
		},
		Params: metahashRequestTxHash{
			Hash: tx,
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url, _ := t.net.TorrentUrl("get-tx")

	resp, err := helperSend(url, reqJson, post)
	if err != nil || resp == nil {
		return nil, err
	}

	var respStruct metahashResponceGetTx

	err = json.Unmarshal(resp, &respStruct)
	if err != nil {
		return nil, err
	}

	return respStruct.Result, nil
}

func (t *metahashNetworkPublicImpV1) Add(addr Address) error {
	if t.net != DevNetwork {
		return &ErrorNetworkUnsupportedMethod{}
	}

	method := fmt.Sprintf("?act=addWallet&p_addr=%s", addr)

	urls, _ := t.net.ProxyUrl(method)

	_, err := helperSend(urls, nil, get)

	return err
}
