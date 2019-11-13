package register

import (
	"proxycenter/pkg/models"
	"sync"
	"unknwon.dev/clog/v2"
)

type Register []func() []*models.IP


func (r *Register) Add(f func()[]*models.IP) {
	*r  = append(*r,f)
}


func (r *Register) Run(ipChan chan<- *models.IP){
	var wg sync.WaitGroup
	for _,f := range *r{
		wg.Add(1)
		go func(f func() []*models.IP) {
			temp := f()
			//log.Println("[run] get into loop")
			for _, v := range temp {
				//log.Println("[run] len of ipChan %v",v)
				ipChan <- v
			}
			wg.Done()
		}(f)
	}
	wg.Wait()
	clog.Info("All getters finished.")
}
