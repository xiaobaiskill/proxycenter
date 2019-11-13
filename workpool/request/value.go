package request

import (
	"encoding/json"
	"regexp"
	"strings"
)

type PostValue struct {
	Action string `json:"action"` // 请求网站
	Method string `json:"method"` // 请求类型 get post delete put...
	Data   string `json:"data"`   // 请求的data数据
	Proxy  bool   // 是否使用代理ip
}

func (p *PostValue) Transform() (h *HttpClient) {
	h = new(HttpClient)
	h.IsProxy = p.Proxy
	switch strings.ToUpper(p.Method) {
	case "POST":
		fallthrough
	case "PUT":
		fallthrough
	case "DELETE":
		h.Method = strings.ToUpper(p.Method)
	default:
		h.Method = "GET"
	}

	reg := regexp.MustCompile("https")
	if reg.MatchString(p.Action) {
		h.IsHttps = true
	}

	h.Url = p.Action
	m := make(map[string]string)
	err := json.Unmarshal([]byte(p.Data), &m)
	if err != nil {
		return
	}
	h.Query = m
	return
}
