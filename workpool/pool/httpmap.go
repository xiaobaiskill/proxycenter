package pool

import (
	"net/http"
	"sync"
)

type HttpMap struct{
	httpclients map[string]*http.Client
	mu sync.Mutex
}


func (h *HttpMap) Store(key string,httpclient *http.Client) {
	h.mu.Lock()
	h.httpclients[key] = httpclient
	h.mu.Unlock()
}

func (h *HttpMap) Load(key string)(*http.Client,bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r,ok := h.httpclients[key]
	return r,ok
}

// 随机一个
func (h *HttpMap) Random()(string,*http.Client,bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.httpclients) > 0 {
		for k,v := range h.httpclients {
			return k,v,true
		}
	}
	return "",nil,false
}

func (h *HttpMap) Delete(key string) {
	h.mu.Lock()
	delete(h.httpclients, key)
	h.mu.Unlock()
}

func (h *HttpMap) Len()int{
	return len(h.httpclients)
}
