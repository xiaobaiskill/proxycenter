package workpool

import (
	"proxycenter/pkg/setting"
	"proxycenter/workpool/pool"
	"proxycenter/workpool/request"
	"strings"
	"time"
	"unknwon.dev/clog/v2"
)


func Run() {
	go func() {
		pool.GetIpToHttp()

		tc := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-tc.C:
				if pool.Httpclients.Len() < int(setting.HttpClientsMax) {
					clog.Info("[work.go] [Run] workpool httpClients 数量 %d 低了100,定时执行程序",pool.Httpclients.Len())
					pool.GetIpToHttp()
				}

			case <-request.GetHttpIpChan:
				if pool.Httpclients.Len() > int(setting.HttpClientsMin){
					continue
				}

				clog.Info("[work.go] [Run] workpool httpClients 数量低了: %d",pool.Httpclients.Len())
				pool.GetIpToHttp()
			}
		}
	}()

	for {
		res := <-request.PostValueChan
		go work(res)
	}
}


func work(value request.PostValue) {
	var res request.Result
	client := request.Transform(value)
	b, proxy, err := client.Request()
	if err != nil {
		clog.Warn("[work.go] [Work] err : %v", err)
		res.CodeStatus = 0
	} else {
		res.CodeStatus = 200
	}

	res.Meta = strings.Trim(string(b), " \n")
	res.Proxy = proxy

	value.Result<-res
}

