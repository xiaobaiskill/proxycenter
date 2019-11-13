package controllers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"proxycenter/workpool"
	"proxycenter/workpool/request"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		param, err := ioutil.ReadAll(r.Body)
		if len(param) == 0 || err != nil {
			io.WriteString(w, "你要做什么？\n")
			return
		}

		//log.Println("传入数据：",string(param))
		value := request.PostValue{}
		json.Unmarshal(param, &value)

		query := r.URL.Query()
		proxy := query["proxy"][0]
		b, _ := strconv.ParseBool(proxy)
		value.Proxy = b

		rc := workpool.Work(value)

		result, ok := <-rc
		if ok {
			b,_ := json.Marshal(result)

			w.Write(b)
		}
	} else {
		io.WriteString(w, "不支持 get 请求")
	}
}
