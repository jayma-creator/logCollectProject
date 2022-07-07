package es

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type LogMessage struct {
	Topic   string
	Message string
}

var esClient *elastic.Client

func Init(addr string) (err error) {
	esClient, err = elastic.NewClient(elastic.SetURL("http://" + addr))
	if err != nil {
		panic(err)
	}
	logrus.Infof("%#v\n", esClient)

	logrus.Info("connect to es success")
	// 从通道中取出数据,写入到es中去
	logrus.Info("Init ES success")
	return
}

func SendToES(topic string, value []byte) {
	msg := &LogMessage{
		Topic:   topic,
		Message: string(value),
	}
	put1, err := esClient.Index().
		Index(topic).
		BodyJson(msg).
		Do(context.Background())
	if err != nil {
		logrus.Error("send to es failed,", err)
	}
	logrus.Infof("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
