package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"time"
)

var cli client.Client
var lastNetIOSTimeStamp int64
var lastNetInfo *NetInfo

func main() {
	err := InitconnInflux()
	if err != nil {
		fmt.Println(err)
		return
	}
	run(time.Second)

}

func InitconnInflux() (err error) {
	cli, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "admin",
		Password: "",
	})
	if err != nil {
		log.Fatal(err)
	}
	return
}

func run(interval time.Duration) {
	ticker := time.Tick(interval)
	for _ = range ticker {
		getCpuInfo()
		getMemInfo()
		getDiskInfo()
		getNetInfo()
	}
}
