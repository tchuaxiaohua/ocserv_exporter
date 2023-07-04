package utils

import (
	"log"
	"strconv"
	"strings"
)

var bandWith float64

func StrToFloat(data string) float64 {
	var (
		data_v string
		data_t string
	)
	dataLst := strings.Split(data, " ")
	if len(dataLst) != 2 {
		data_v = ""
		data_t = ""
	} else {
		data_v = dataLst[0]
		data_t = dataLst[1]
	}

	switch data_t {
	case "bytes/sec":
		v, err := strconv.ParseFloat(data_v, 64)
		if err != nil {
			log.Printf("数据转换出错：%s\n", err)
		}
		bandWith = v / 1024
	case "KB/sec":
		v, err := strconv.ParseFloat(data_v, 64)
		if err != nil {
			log.Printf("数据转换出错：%s\n", err)
		}
		bandWith = v
	case "MB/sec":
		v, err := strconv.ParseFloat(data_v, 64)
		if err != nil {
			log.Printf("数据转换出错：%s\n", err)
		}
		bandWith = v * 1024
	default:
		bandWith = 0
	}
	return bandWith
}
