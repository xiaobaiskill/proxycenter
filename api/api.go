package api

import (
	"fmt"
	"net/http"
	"proxycenter/api/controllers"
	"proxycenter/api/middleware"
	"proxycenter/pkg/setting"
)

func Run() {
	// 初识化路由
	route := router()

	// 初识化中间件
	m := middleWare.NewMiddleWare(route)

	// http
	addr := fmt.Sprintf("%s:%s", setting.AppAddr, setting.AppPort)
	http.ListenAndServe(addr, m)
}

// 路由
func router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/index", controllers.Index)
	return mux
}
