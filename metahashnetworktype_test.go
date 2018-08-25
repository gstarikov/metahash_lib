package metahash_lib

import (
	"testing"
)

func TestMetahashNetworkType(t *testing.T) {
	nt := NetworkType(DevNetwork)
	_, err := nt.ProxyUrl("")
	if err != nil {
		t.Errorf("%v", err)
	}

	_, err = nt.TorrentUrl("any_method")
	if err != nil {
		t.Errorf("%v", err)
	}
}
