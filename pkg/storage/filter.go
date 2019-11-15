package storage

import (
	"proxycenter/pkg/models"
	"proxycenter/proxypool/check/register"
	clog "unknwon.dev/clog/v2"

	"sync"
)

// CheckProxy .
func CheckProxy(ip *models.IP) {
	if CheckIP(ip) {
		ProxyAdd(ip)
	}
}

// CheckIP is to check the ip work or not
func CheckIP(ip *models.IP) bool {
	if register.Check.Run(ip) {
		if err := models.Update(*ip); err != nil {
			clog.Warn("[CheckIP] Update IP = %v Error = %v", *ip, err)
		}
		return true
	}
	return false
}

// CheckProxyDB to check the ip in DB
func CheckProxyDB() {
	x := models.CountIPs()
	clog.Info("Before check, DB has: %d records.", x)
	ips, err := models.GetAll()
	if err != nil {
		clog.Warn(err.Error())
		return
	}
	var wg sync.WaitGroup
	for _, v := range ips {
		wg.Add(1)
		go func(v *models.IP) {
			if !CheckIP(v) {
				ProxyDel(v)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	x = models.CountIPs()
	clog.Info("[filter.go] [CheckProxyDB] After check, DB has: %d records.", x)
}

// ProxyRandom .
func ProxyRandom() (ip *models.IP) {
	ips, err := models.GetAll()
	x := len(ips)
	clog.Info("[fileter.go][ProxtRandom] http: len(ips) = %d", x)
	if err != nil || x == 0 {
		clog.Warn("[fileter.go][ProxtRandom]: 没有获取到ip")
		return models.NewIP()
	}
	randomNum := RandInt(0, x)

	return ips[randomNum]
}

func ProxyFindsToNum(num int)(ips []*models.IP, ok bool){
	ips, err := models.GetAllToNum(num)
	if err != nil {
		clog.Warn("[fileter.go] [ProxyFindsToNum] 获取 代理ip 失败")
		return
	}
	ok = true
	return
}


// ProxyFind .
func ProxyFind(value string) (ip *models.IP) {
	ips, err := models.FindAll(value)
	if err != nil {
		clog.Warn("[fileter.go] [ProxyFind] 获取 %v 代理ip 失败",value)
		return models.NewIP()
	}
	x := len(ips)
	clog.Warn("[fileter.go] [ProxyFind] %s: len(ips) = %d",value, x)
	randomNum := RandInt(0, x)
	clog.Info("[fileter.go] [proxyFind] random num = %d", randomNum)
	if randomNum == 0 {
		return models.NewIP()
	}
	return ips[randomNum]
}

// ProxyAdd .
func ProxyAdd(ip *models.IP) {
	models.InsertIps(ip)
}

// ProxyDel .
func ProxyDel(ip *models.IP) {
	models.DeleteIP(ip)
}
