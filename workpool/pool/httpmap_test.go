package pool

import (
	"net/http"
	"testing"
)

func TestHttpMap(t *testing.T) {
	var hm HttpMap
	hm.httpclients = make(map[string]*http.Client)
	hm.Store("1111",&http.Client{Timeout:1111})
	hm.Store("2222",&http.Client{Timeout:2222})
	hm.Store("4444",&http.Client{Timeout:4444})
	hm.Store("5555",&http.Client{Timeout:5555})
	hm.Store("6666",&http.Client{Timeout:6666})

	v,ok := hm.Load("6666")
	if ok && v.Timeout == 6666{
		t.Log("httpmap 读取成功",v)
	} else {
		t.Fatal("httpmap 读取成功")
	}

	hm.Delete("4444")

	v,ok = hm.Random()
	if !ok {
		t.Fatal("随机读取失败")
	}
}
