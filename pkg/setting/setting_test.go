package setting

import "testing"

func TestNewContext(t *testing.T) {
	i, err := NewContext("../../conf/app.ini")
	if err != nil {
		t.Fatal("文件加载失败")
	}

	if i.GetString("", "APP_NAME", "") != "proxycenter" {
		t.Log("APP_NAME != proxycenter")
	}
}
