package register

import (
	"proxycenter/pkg/models"
	"proxycenter/proxypool/check"
)

type Register []func(*models.IP)bool


func (r *Register) Add(f func(*models.IP)bool) {
	*r  = append(*r,f)
}


func (r *Register) Run(ip *models.IP)bool{
	for _,f := range *r{
		b := f(ip)
		if !b{
			return false
		}
	}
	return true
}



var Check Register
func init(){
	Check.Add(check.Mlytics)
	//Check.Add(check.Httpbin)
}
