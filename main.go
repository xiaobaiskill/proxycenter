package main

import (
	"proxycenter/api"
	"proxycenter/pkg/initial"
	"proxycenter/proxypool"
	"proxycenter/workpool"
)

func main() {
	initial.GlobalInit()
	go proxypool.Run()
	go workpool.Run()
	api.Run()
}
