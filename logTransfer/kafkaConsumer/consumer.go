package kafkaConsumer

import (
	"ch06-qimiProject/logTransfer/es"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

func Init(addr []string, topic string) (err error) {
	// 创建新的消费者
	consumer, err := sarama.NewConsumer(addr, nil)
	if err != nil {
		logrus.Errorf("fail to start consumer, err:%v\n", err)
		return
	}

	// 拿到指定topic下面的所有分区列表
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		logrus.Errorf("fail to get list of partition:err%v\n", err)
		return
	}
	logrus.Info(partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		var pc sarama.PartitionConsumer
		pc, err = consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			logrus.Errorf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		// 异步从每个分区消费信息
		logrus.Info("start to consume...")
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				logrus.Info(msg.Topic, string(msg.Value))
				var m1 map[string]interface{}
				err = json.Unmarshal(msg.Value, &m1)
				if err != nil {
					logrus.Errorf("unmarshal msg failed, err:%v\n", err)
					continue
				}
				// 为了将同步流程异步化,所以将取出的日志数据先放到channel中
				es.PutLogData(m1)
			}
		}(pc)
	}
	logrus.Info("Init kafka success")
	return
}
