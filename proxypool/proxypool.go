package proxypool

import (
	"proxycenter/pkg/models"
	"proxycenter/pkg/storage"
	"proxycenter/proxypool/getter/register"
	"time"
	clog "unknwon.dev/clog/v2"
)

func Run() {
	ipChan := make(chan *models.IP, 2000)

	// Check the IPs in DB
	go func() {
		storage.CheckProxyDB()
	}()

	for i := 0; i < 50; i++ {
		go func() {
			for {
				storage.CheckProxy(<-ipChan)
			}
		}()
	}

	// Start getters to scraper IP and put it in channel
	for {
		n := models.CountIPs()
		clog.Info("[proxycenter.go] Chan: %v, IP: %v\n", len(ipChan), n)
		if len(ipChan) < 100 {
			go register.Reg.Run(ipChan)
		}
		time.Sleep(5 * time.Minute)
	}
}
