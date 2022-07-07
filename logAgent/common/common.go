package common

import (
	"fmt"
	"net"
	"strings"
)

type Config struct {
	*KafkaConfig  `ini:"kafka"`
	*CollectEntry `ini:"collect"`
	*EtcdConfig   `ini:"etcd"`
}
type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int64  `ini:"chan_size"`
}

type EtcdConfig struct {
	Address    []string `ini:"address"`
	CollectKey string   `ini:"collect_key"`
}

//要收集的日志的配置项结构体
type CollectEntry struct {
	Path  string `json:"path"`  //去哪个路径读取日志文件
	Topic string `json:"topic"` //日志文件发往kafka的哪个topic
}

func GetIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer conn.Close()
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
	return ip
}
