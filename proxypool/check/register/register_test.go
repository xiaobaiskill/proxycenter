package register

import (
	"proxycenter/pkg/models"
	"testing"
)

func TestRegister(t *testing.T) {
	ip := new(models.IP)
	ip.Data = "157.245.88.191:8080"  // 清输入可用的ip代理
	ip.Type1 = "http"

	if !Check.Run(ip) {
		t.Log("[check][register] 代理添加失败")
	}
}
