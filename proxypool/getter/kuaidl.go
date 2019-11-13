package getter

import (
	"github.com/Aiicy/htmlquery"
	"proxycenter/pkg/models"
	clog "unknwon.dev/clog/v2"
)

// KDL get ip from kuaidaili.com
func KDL() (result []*models.IP) {
	pollURL := "http://www.kuaidaili.com/free/inha/"
	doc, err := htmlquery.LoadURL(pollURL)
	if err !=nil {
		clog.Error("[kuaidl] 数据爬取失败：%+v",err)
		return
	}
	trNode, err := htmlquery.Find(doc, "//table[@class='table.table-bordered.table-striped']//tbody//tr")
	if err != nil {
		clog.Warn("[kuaidl] 解析失败：%+v" ,err.Error())
		return
	}
	for i := 0; i < len(trNode); i++ {
		tdNode, _ := htmlquery.Find(trNode[i], "//td")
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		Type := htmlquery.InnerText(tdNode[3])
		speed := htmlquery.InnerText(tdNode[5])

		IP := models.NewIP()
		IP.Data = ip + ":" + port
		if Type == "HTTPS" {
			IP.Type1 = ""
			IP.Type2 = "https"
		} else if Type == "HTTP" {
			IP.Type1 = "http"
		}
		IP.Speed = extractSpeed(speed)
		result = append(result, IP)
	}

	clog.Info("[kuaidaili] done")
	return
}
