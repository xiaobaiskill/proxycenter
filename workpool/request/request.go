package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"proxycenter/pkg/models"
	"proxycenter/pkg/setting"
	"proxycenter/pkg/storage"
	"proxycenter/workpool/pool"
	"regexp"
	"strings"
	"time"
	"unknwon.dev/clog/v2"
)

// 返回结果
type Result struct {
	Meta       string `json:"meta"`
	Proxy      bool   `json:"proxy"`
	CodeStatus int    `json:"code"` // 返回状态 200 正常    0 错误
}

// 请求的参数
type PostValue struct {
	Action string      `json:"action"` // 请求网站
	Method string      `json:"method"` // 请求类型 get post delete put...
	Data   string      `json:"data"`   // 请求的data数据
	Proxy  bool        // 是否使用代理ip
	Result chan Result // 返回数据使用的chan
}

type HttpClient struct {
	Method  string
	IsProxy bool
	Url     string
	Query   map[string]string
	IsHttps bool
	IP      *models.IP
}

func (h *HttpClient) Request() (body []byte, proxy bool, err error) {
	req, err := http.NewRequest(h.Method, h.Url, nil)
	if err != nil {
		return
	}
	if len(h.Query) != 0 {
		q := req.URL.Query()
		for k, v := range h.Query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	var (
		resp   *http.Response
		client *http.Client
		ip     string
		ok     bool
	)
	proxy = true
	if h.IsProxy {
		ip, client, ok = RandClient()
		if !ok {
			clog.Error("[request.go] [Request]没有从workpool 中获取到httpclient")
			return
		}
		clog.Info("[request.go] [Request]使用的代理ip :%s", ip)
	} else {
		proxy = false
		client = &http.Client{Timeout: time.Duration(setting.TimeOut) * 1000 * 1000}
	}


	resp, err = client.Do(req)
	if err != nil {
		storage.ProxyDel(&models.IP{Data: ip})
		pool.Httpclients.Delete(ip)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

var (
	GetHttpIpChan = make(chan bool) //httpclient 不足时发送指令从新获取
	PostValueChan = make(chan PostValue, 100)
)

// 随机获取
func RandClient() (ip string, client *http.Client, ok bool) {
	if pool.Httpclients.Len() < 10 && len(GetHttpIpChan) == 0 {
		GetHttpIpChan <- true
	}
	ip, client, ok = pool.Httpclients.Random()
	return

}

// post 请求数据 转 httpclient
func Transform(p PostValue) (h *HttpClient) {
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
