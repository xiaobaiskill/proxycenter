package workpool

import (
	"proxycenter/workpool/request"
	"strings"
	"unknwon.dev/clog/v2"
)

type Result struct {
	Meta  string `json:"meta"`
	Proxy bool   `json:"proxy"`
	CodeStatus   int  `json:"code"`  // 返回状态 200 正常    0 错误
}

func Work(value request.PostValue) (rc chan Result) {
	rc = make(chan Result)
	go func() {
		var res Result
		b, proxy, err := (&value).Transform().Request()
		if err !=nil {
			clog.Warn("[work.go] [Work] err : %v",err)
			res.CodeStatus = 0
		} else {
			res.CodeStatus = 200
		}

		res.Meta = strings.Trim(string(b)," \n")
		res.Proxy = proxy
		rc <- res
	}()
	return rc
}
