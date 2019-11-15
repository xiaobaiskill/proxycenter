package getter

import (
	clog "unknwon.dev/clog/v2"

	"github.com/Aiicy/htmlquery"
	"proxycenter/pkg/models"
	"regexp"
	"strconv"
)

//feiyi get ip from feiyiproxy.com
func Feiyi() (result []*models.IP) {
	clog.Info("[FEIYI] start test")
	pollURL := "http://www.feiyiproxy.com/?page_id=1457"
	doc, _ := htmlquery.LoadURL(pollURL)
	trNode, err := htmlquery.Find(doc, "//div[@class='et_pb_code et_pb_module  et_pb_code_1']/div/table/tbody/tr")
	clog.Info("[FEIYI] start up")
	if err != nil {
		clog.Info("FEIYI] parse pollUrl error")
		clog.Warn(err.Error())
	}
	//debug begin
	clog.Info("[FEIYI] len(trNode) = %d ", len(trNode))
	for i := 1; i < len(trNode); i++ {
		tdNode, _ := htmlquery.Find(trNode[i], "//td")
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		Type := htmlquery.InnerText(tdNode[3])
		speed := htmlquery.InnerText(tdNode[6])

		IP := models.NewIP()
		IP.Data = ip + ":" + port

		if Type == "HTTPS" {
			IP.Type1 = "https"
			IP.Type2 = ""

		} else if Type == "HTTP" {
			IP.Type1 = "http"
		}
		IP.Speed = pyextractSpeed(speed)

		clog.Info("[FEIYI] ip.Data = %s,ip.Type = %s,%s ip.Speed = %d", IP.Data, IP.Type1, IP.Type2, IP.Speed)

		result = append(result, IP)
	}

	clog.Info("FEIYI done.")
	return
}

func extractSpeed(oritext string) int64 {
	reg := regexp.MustCompile(`\[1-9\]\d\*\\.\?\d\*`)
	temp := reg.FindString(oritext)
	if temp != "" {
		speed, _ := strconv.ParseInt(temp, 10, 64)
		return speed
	}
	return -1
}

// 支持 3秒 0.3秒这中的形式
func pyextractSpeed(oritext string) int64 {
	reg := regexp.MustCompile(`\d+?\.?\d*`)
	temp := reg.FindStringSubmatch(oritext)

	if len(temp) >= 1 && temp[0] != "" {
		speed_float, _ := strconv.ParseFloat(temp[0],64)
		speed := int64(speed_float * 1000)
		return speed
	}
	return -1
}
