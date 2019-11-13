package getter

import (
	"testing"
)


func TestPyextractSpeed(t *testing.T) {
	text := "31.322314秒"
	i := pyextractSpeed(text)

	if i != 31 {
		t.Fatal("[TestPyextractSpeed] 数据解析失败")
	}
}
