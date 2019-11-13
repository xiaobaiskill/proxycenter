package request

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"proxycenter/pkg/models"
	"proxycenter/pkg/setting"
	"proxycenter/pkg/storage"
	"time"
	clog "unknwon.dev/clog/v2"
)

type HttpClient struct {
	Method  string
	IsProxy bool
	Url     string
	Query   map[string]string
	IsHttps bool
}

func (h *HttpClient) ProxyHttpClient() (*http.Client, bool) {
	if h.IsProxy {
		var ip *models.IP
		//if h.IsHttps {
		//	ip = storage.ProxyFind("https")
		//} else {
			ip = storage.ProxyRandom()
		//}
		
		clog.Info("[request] [ProxyHttpClient] ip data:%v",ip)
		if ip.Data == "" {
			clog.Info("[request] [ProxyHttpClient] 没有获取到可代理的ip")
			return &http.Client{
				Timeout:time.Duration(setting.TimeOut) * 1000 * 1000,
			}, false
		}
		proxy := func(_ *http.Request) (*url.URL, error) {
			http := "http://"
			//if h.IsHttps {
			//	http = "https://"
			//}
			return url.Parse(http + ip.Data)
		}

		transport := &http.Transport{Proxy: proxy}
		client := &http.Client{
			Transport: transport,
			Timeout:time.Duration(setting.TimeOut) * 1000 * 1000,  //  1 second == 1000 000 000 time.Duration
		}

		return client, true
	} else {
		return &http.Client{
			Timeout:time.Duration(setting.TimeOut) * 1000 * 1000,
		}, false
	}
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

	var resp *http.Response
	var client *http.Client
	client, proxy = h.ProxyHttpClient()
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}
