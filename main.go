package main

import (
	"proxycenter/api"
	"proxycenter/pkg/initial"
	"proxycenter/proxypool"
)

func main() {
	initial.GlobalInit()
	go proxypool.Run()
	api.Run()
}
