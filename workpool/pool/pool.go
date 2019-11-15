package pool

import (
	"net/http"
	"net/url"
	"proxycenter/pkg/models"
	"proxycenter/pkg/setting"
	"proxycenter/pkg/storage"
	"time"
	"unknwon.dev/clog/v2"
)

var (
	Httpclients  HttpMap // http
)

func init(){
	Httpclients.httpclients = make(map[string]*http.Client)
}

// 获取ip
func GetIpToHttp() {
	clog.Info("[pool.go] [GetIpToHttp]获取代理ip数据")
	ips, ok := storage.ProxyFindsToNum(int(setting.HttpClientLimit))

	if ok {
		for _, v := range ips {
			Httpclients.Store(v.Data,createHttpClient(v))
		}
	}
	clog.Info("[pool.go] [GetIpToHttp] workpool 池的数量:%d",Httpclients.Len())
}

// 通过IP 生成 httpclient
func createHttpClient(ip *models.IP) *http.Client {
	proxy := func(_ *http.Request) (*url.URL, error) {
		http := "http://"

		if ip.Type2 == "https" {
			http = "https://"
		}
		return url.Parse(http + ip.Data)
	}

	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(setting.TimeOut) * 1000 * 1000, //  1 second == 1000 000 000 time.Duration
	}

	return client
}
