package kafkaProducer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var Client sarama.SyncProducer
var msgChan chan *sarama.ProducerMessage

func Init(address []string, chanSize int64) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	//初始化连接
	Client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("kafka: producer closed, err:", err)
		return
	}

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
			fmt.Println(pid, offset)
		}
	}
}

//对外暴露msg，接收信息
func ToMsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
