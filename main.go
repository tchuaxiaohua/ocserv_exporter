package main

import (
	"fmt"
	"github.com/tchuaxiaohua/ocserv_exporter/cmd"
)

func main() {
	if err := cmd.StartCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
