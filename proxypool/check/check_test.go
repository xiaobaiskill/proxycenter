package check

import (
	"proxycenter/pkg/models"
	"testing"
)

func TestMlytics(t *testing.T){
	ip := new(models.IP)
	ip.Data = "157.245.88.191:8080"  // 清输入可用的ip代理
	ip.Type1 = "http"

	b := Mlytics(ip)
	if !b {
		t.Log("免费ip代理,测试失败")
	}
}