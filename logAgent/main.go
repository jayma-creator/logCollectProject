package main

import (
	"ch06-qimiProject/logAgent/common"
	"ch06-qimiProject/logAgent/etcd"
	"ch06-qimiProject/logAgent/kafkaProducer"
	"ch06-qimiProject/logAgent/tailfile"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

var configObj = new(common.Config)

func run() {
	select {}
}

func main() {
	ip := common.GetIP()
	file := "./logAgent/conf/logAgent.ini"
	err := ini.MapTo(configObj, file)
	if err != nil {
		logrus.Errorf("load config failed, err:%v", err)
		return
	}
	logrus.Info("load config success")

	//连接kafka，循环读取消息管道
	err = kafkaProducer.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed, err:%v", err)
		return
	}

	//连接etcd
	err = etcd.Init(configObj.EtcdConfig.Address)
	if err != nil {
		logrus.Errorf("init etcd failed, err:%v", err)
		return
	}
	//根据key拉取要收集的日志配置，不同ip对应不同的key
	collectKey := fmt.Sprintf(configObj.EtcdConfig.CollectKey, ip)
	allConf, err := etcd.GetConf(collectKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	//监控etcd的key对应值的变化
	go etcd.WatchConf(collectKey)

	//连接tail，把日志信息发往消息管道
	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Errorf("init tailfile failed,err:%v", err)
		return
	}
	run()
}
