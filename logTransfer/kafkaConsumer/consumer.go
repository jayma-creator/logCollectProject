package kafkaConsumer

import (
	"ch06-qimiProject/logTransfer/es"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"sync"
)

var wg sync.WaitGroup

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
		defer pc.AsyncClose()
		// 异步从每个分区消费
		logrus.Info("start to consume...")
		go func(sarama.PartitionConsumer) {
			wg.Add(1)
			defer wg.Done()
			for msg := range pc.Messages() {
				logrus.Info(msg.Topic, string(msg.Value))
				es.SendToES(msg.Topic, msg.Value)
				if err != nil {
					logrus.Errorf("unmarshal msg failed, err:%v\n", err)
					continue
				}
			}
		}(pc)
	}
	wg.Wait()
	logrus.Info("Init kafka success")
	return
}
