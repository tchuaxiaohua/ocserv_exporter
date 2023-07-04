package ocserv

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"log"

	"github.com/tchuaxiaohua/ocserv_exporter/utils"
)

var (
	clients float64
)

type OcservUsersCollector struct {
	desc  *prometheus.Desc
	Hosts []string
	Port  string
	Team  string
}

// NewOcservUsersCollector Docker及远程主机模式
func NewOcservUsersCollector(ip []string, port string) *OcservUsersCollector {
	res := &OcservUsersCollector{
		desc: prometheus.NewDesc(
			"ocserv_client_count",
			"ocserv 客户端链接数",
			[]string{"instance", "team"},
			nil),
		Hosts: ip,
		Port:  port,
		Team:  "ops",
	}
	return res
}

type Data struct {
	User string `json:"Username"`
}

func (p *OcservUsersCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- p.desc
}

func (p *OcservUsersCollector) Collect(metrics chan<- prometheus.Metric) {
	for _, host := range p.Hosts {
		lst := []Data{}
		log.Printf("开始采集主机%s指标\n", host)
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
		clients = float64(len(lst))
		metrics <- prometheus.MustNewConstMetric(p.desc, prometheus.CounterValue, clients, host, p.Team)
		log.Printf("主机%s指标采集完毕\n", host)
	}
}
