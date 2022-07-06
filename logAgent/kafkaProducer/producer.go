package kafkaProducer

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var Client sarama.SyncProducer
var msgChan chan *sarama.ProducerMessage

func Init(address []string, chanSize int64) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //新选出一个partition
	config.Producer.Return.Successes = true                   //成功交付的消息在success channel返回

	//初始化连接
	Client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("连接Kafka失败:", err)
		return
	}
	defer Client.Close()
	//信息管道
	msgChan = make(chan *sarama.ProducerMessage, chanSize)
	go sendMsg()
	logrus.Info("init kafka success")
	return
}

//kafka循环读取管道
func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := Client.SendMessage(msg)
			if err != nil {
				logrus.Error("send msg failed, err:", err)
				return
			}
			logrus.Info("send msg to kafka success.")
			logrus.Info(pid, offset)
		}
	}
}

//对外暴露msg，接收信息
func ToMsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
