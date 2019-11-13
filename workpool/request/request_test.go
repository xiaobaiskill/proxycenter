package request

import (
	"encoding/json"
	"testing"
)

func TestPostValue_Transform(t *testing.T) {
	p := PostValue{
		"https://39.103.32.1",
		"get",
		`{"name":"jmz","age":"12"}`,
		true,
	}
	hc := p.Transform()
	if hc.IsProxy == true && hc.Query["age"] == "12" && hc.Method == "GET" && hc.Url == "https://39.103.32.1" && hc.IsHttps{
		t.Log("转换正确")
	} else {
		t.Fatal("转换失败")
	}
}


func TestSubmit(t *testing.T) {
	type Finfo struct {
		Id            string `json:"id"`
		Title         string `json:"title"`
		Creattime     string `json:"creattime"`
		Updatetime    string `json:"updatetime"`    // updatetime 可能是用户最近的更新的时间
		Lucupdatetime string `json:"lucupdatetime"` //  365 房产用的是这个时间为发布时间 ：lucupdatetime
		Contactor     string `json:"contactor"`
		Telno         string `json:"telno"`
	}

	p := PostValue{
		"https://m.house365.com/index.php",
		"post",
		`{"m":"home","c":"Secondhouse","a":"index","q":"i0_j0_x0_m0_h0_f0_z0_o0_k","act":"search","page":"2","city":"nj"}`,
		false,
	}
	client := p.Transform()

	b, _, err := client.Request()
	if err != nil {
		t.Fatal("请求失败：" + err.Error())
	}
	info := make([]Finfo, 1)
	err = json.Unmarshal(b, &info)
	if err != nil {
		t.Fatal("内容解析失败：" + err.Error())
	}
	t.Log(info)
}
