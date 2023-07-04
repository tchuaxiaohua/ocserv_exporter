package ocserv

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	flag float64
)

// 实现 prometheus Collector接口
type OcservCollector struct {
	desc  *prometheus.Desc
	Hosts []string
	Port  string
}

// NewOcservCollector Docker及远程主机模式
func NewOcservCollector(ip []string, port string) *OcservCollector {
	res := &OcservCollector{
		desc: prometheus.NewDesc(
			"ocserv_status",
			fmt.Sprintf("ocserv 端口 %s 存活状态", port),
			[]string{"hostname"},
			nil),
		Hosts: ip,
		Port:  port,
	}
	return res
}

func (p *OcservCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- p.desc
}

func (p *OcservCollector) Collect(metrics chan<- prometheus.Metric) {
	for _, host := range p.Hosts {
		cmd := exec.Command("ncat", "-vz", host, p.Port)
		out, _ := cmd.CombinedOutput()
		if strings.Contains(string(out), "bytes sent") || strings.Contains(string(out), "open") {
			flag = 0
		} else {
			flag = 1
		}
		metrics <- prometheus.MustNewConstMetric(p.desc, prometheus.CounterValue, flag, host)
	}
}
