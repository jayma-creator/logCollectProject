package main

import (
	"ch06-qimiProject/logTransfer/common"
	"ch06-qimiProject/logTransfer/es"
	"ch06-qimiProject/logTransfer/kafkaConsumer"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

func main() {
	//加载配置文件
	var cfg = new(common.Config)
	err := ini.MapTo(cfg, "./logTransfer/config/logTransfer.ini")
	if err != nil {
		logrus.Errorf("load config failed,err:%v\n", err)
		return
	}
	logrus.Info("load config success")

	//连接ES
	err = es.Init(cfg.ESConf.Address)
	if err != nil {
		logrus.Errorf("Init es failed,err:%v\n", err)
		return
	}

	//连接kafka
	err = kafkaConsumer.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.Topic)
	if err != nil {
		logrus.Errorf("connect to kafka failed,err:%v\n", err)
		return
	}
	// 在这儿停顿
	select {}
}
