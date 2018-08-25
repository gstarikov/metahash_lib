package metahash_lib

import (
	"net"
)

type NetworkType int

const (
	DevNetwork  NetworkType = iota
	ProdNetwork             = iota
)

type networkSubType int

const (
	tor   networkSubType = iota
	proxy                = iota
)

func (t NetworkType) host() string {
	switch t {
	case DevNetwork:
		return "net-dev.metahashnetwork.com"
	case ProdNetwork:
		return "net-main.metahashnetwork.com"
	default:
		panic("undefined behavior")
	}
}

type NetworkErrorCantResolve struct{}

func (t *NetworkErrorCantResolve) Error() string {
	return "NetworkErrorCantResolve"
}

func (t NetworkType) address(sNet networkSubType) ([]string, error) {
	host := t.host()
	switch sNet {
	case tor:
		host = "tor." + host
	case proxy:
		host = "proxy." + host
	default:
		return nil, &NetworkErrorCantResolve{}
	}
	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil, err
	}
	if len(addrs) > 0 {
		return addrs, nil
	}
	return nil, &NetworkErrorCantResolve{}
}

func (t NetworkType) ProxyUrl(method string) ([]string, error) {
	addr, err := t.address(proxy)
	if err != nil {
		return nil, err
	}
	for i := range addr {
		addr[i] = "http://" + addr[i] + ":9999/" + method
	}
	return addr, nil
}

func (t NetworkType) TorrentUrl(method string) ([]string, error) {
	addr, err := t.address(tor)
	if err != nil {
		return nil, err
	}
	for i := range addr {
		addr[i] = "http://" + addr[i] + ":5795/" + method
	}
	return addr, nil
}
