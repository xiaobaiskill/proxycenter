package controllers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
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

		resultChan := make(chan request.Result)

		value.Result = resultChan

		request.PostValueChan <- value

		result, ok := <-resultChan
		if ok {
			b,_ := json.Marshal(result)

			w.Write(b)
		}
	} else {
		io.WriteString(w, "不支持 除post之外的其他请求")
	}
}
