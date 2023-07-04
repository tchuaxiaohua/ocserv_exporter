package ocserv

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/tchuaxiaohua/ocserv_exporter/utils"
)

type OcservUsersDetailCollector struct {
	desc  *prometheus.Desc
	Hosts []string
	Port  string
	Team  string
}

type UsersData struct {
	Id        int    `json:"ID"`
	User      string `json:"Username"`
	AverageRX string `json:"Average RX"`
	AverageTX string `json:"Average TX"`
}

func NewOcservUsersDetailCollector(ip []string, port string) *OcservUsersDetailCollector {
	res := &OcservUsersDetailCollector{
		desc: prometheus.NewDesc(
			"ocserv_client_info",
			"ocserv 客户端详情",
			[]string{"instance", "id", "bandwidth", "username"},
			nil),
		Hosts: ip,
		Port:  port,
		Team:  "ops",
	}
	return res
}

func (p *OcservUsersDetailCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- p.desc
}

func (p *OcservUsersDetailCollector) Collect(metrics chan<- prometheus.Metric) {
	for _, host := range p.Hosts {
		lst := []UsersData{}
		log.Printf("开始采集主机【%s】 ocserv 客户端详情指标\n", host)
		cmds, err := utils.Connect(host, p.Port)
		if err != nil {
			log.Printf("主机连接异常:%s\n", err)
			continue
		}
		err = json.Unmarshal(cmds, &lst)
		if err != nil {
			log.Printf("数据解析失败:%s\n", err)
			continue
		}

		for _, item := range lst {
			metrics <- prometheus.MustNewConstMetric(p.desc, prometheus.CounterValue, utils.StrToFloat(item.AverageRX), host, strconv.Itoa(item.Id), "send", item.User)
			metrics <- prometheus.MustNewConstMetric(p.desc, prometheus.CounterValue, utils.StrToFloat(item.AverageTX), host, strconv.Itoa(item.Id), "receive", item.User)
		}
	}
}
