package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tchuaxiaohua/ocserv_exporter/ocserv"
	"github.com/tchuaxiaohua/ocserv_exporter/utils"
)

var (
	ocservHost string
	hostPort   string
	ocservPort string
	confFile   string
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "启动服务",
	Long:  "启动服务",
	RunE: func(cmd *cobra.Command, args []string) error {

		ocservList := []string{ocservHost}
		prometheus.MustRegister(ocserv.NewOcservCollector(ocservList, ocservPort))
		prometheus.MustRegister(ocserv.NewOcservUsersCollector(ocservList, hostPort))
		prometheus.MustRegister(ocserv.NewOcservUsersDetailCollector(ocservList, hostPort))

		//	暴露指标
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe("0.0.0.0:18086", nil)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

func loadConfig() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("etc/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed: %v", err)
	}
}

func init() {
	// 加载配置文件
	loadConfig()
	// 配置
	StartCmd.PersistentFlags().StringVarP(&utils.SshKeyPath, "key", "k", viper.GetString("host.sshKey"), "主机远程秘钥文件路径")
	StartCmd.PersistentFlags().StringVarP(&utils.SshPassword, "pwd", "a", viper.GetString("host.passWord"), "主机远程密码")
	StartCmd.PersistentFlags().StringVarP(&ocservHost, "ip", "i", viper.GetString("host.ip"), "ocserv主机远IP")
	StartCmd.PersistentFlags().StringVarP(&hostPort, "port", "P", viper.GetString("host.port"), "ocserv主机端口")
	StartCmd.PersistentFlags().StringVarP(&ocservPort, "svcport", "p", viper.GetString("app.port"), "ocserv服务端口")
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/app.yaml", "项目配置文件路径")
}
