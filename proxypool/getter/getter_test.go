package getter

import (
	"testing"
	"unknwon.dev/clog/v2"
)


func TestPyextractSpeed(t *testing.T) {
	text := "0.4秒"
	i := pyextractSpeed(text)
	if i != 400 {
		t.Fatal("[TestPyextractSpeed] 数据解析失败")
	}
}

func TestKDL(t *testing.T) {
	clog.NewConsole()

	data := KDL()
	if len(data) == 0 {
		t.Log("KDL 没有获取到数据")
	}
}



func TestPydl(t *testing.T) {
	clog.NewConsole()

	data := Pydl()
	if len(data) == 0 {
		t.Log("Pyl 没有获取到数据")
	}
}


func TestFeiyi(t *testing.T) {
	clog.NewConsole()

	data := Feiyi()
	if len(data) == 0 {
		t.Log("Feiyi 没有获取到数据")
	}
}

func TestIP66(t *testing.T) {
	clog.NewConsole()

	data := IP66()
	if len(data) == 0 {
		t.Log("IP66 没有获取到数据")
	}
}

func TestIP89(t *testing.T) {
	clog.NewConsole()

	data := IP89()
	if len(data) == 0 {
		t.Log("IP89 没有获取到数据")
	}
}
