package check

import (
	sj "github.com/bitly/go-simplejson"
	"github.com/parnurzeal/gorequest"
	"proxycenter/pkg/models"
	"time"
	"unknwon.dev/clog/v2"
)

//IP89 get ip from www.89ip.cn
func Mlytics(ip *models.IP) (b bool) {
	var pollURL = "https://jimqaweb.mlytics.ai/cache.txt"
	var testIP string
	if ip.Type2 == "https" {
		testIP = "https://" + ip.Data
	} else {
		testIP = "http://" + ip.Data
	}
	begin := time.Now()
	resp, _, errs := gorequest.New().Proxy(testIP).Get(pollURL).End()
	if errs != nil {
		clog.Info("[CheckIP] testIP = %s, pollURL = %s: Error = %v", testIP, pollURL, errs)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		//harrybi 20180815 判断返回的数据格式合法性
		_, err := sj.NewFromReader(resp.Body)
		if err != nil {
			clog.Info("[CheckIP] testIP = %s, pollURL = %s: Error = %v", testIP, pollURL, err)
			return false
		}
		//harrybi 计算该代理的速度，单位毫秒
		ip.Speed = time.Now().Sub(begin).Nanoseconds() / 1000 / 1000 //ms
		return true
	}
	return false

}
