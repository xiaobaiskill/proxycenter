package register

import (
	"proxycenter/pkg/models"
	"testing"
)

func TestRegister(t *testing.T) {
	var reg Register
	reg.Add(func() []*models.IP {
		return []*models.IP{&models.IP{ID:1}}
	})

	ipchan := make(chan *models.IP,1)
	go reg.Run(ipchan)
	v := <- ipchan
	close(ipchan)
	if v.ID != 1 {
		t.Fatal("getter_register 测试未通过")
	}
}
