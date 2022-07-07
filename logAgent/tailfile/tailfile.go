package tailfile

import (
	"ch06-qimiProject/logAgent/kafkaProducer"
	"context"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type tailTask struct {
	path   string
	topic  string
	tObj   *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

func newTailTask(path, topic string) *tailTask {
	ctx, cancel := context.WithCancel(context.Background())
	tt := &tailTask{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
	return tt
}

func (t *tailTask) Init() (err error) {
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //2表示末尾
		MustExist: false,                                //允许不存在
		Poll:      true,                                 //轮询
	}
	//用tail打开文件
	t.tObj, err = tail.TailFile(t.path, cfg)
	return
}

func (t *tailTask) run() {
	//读取日志，发往kafka
	logrus.Infof("collect for path:%s is running...", t.path)
	for {
		select {
		//调用cancel()就会传一个空结构体，结束这个goroutine
		case <-t.ctx.Done():
			logrus.Infof("path:%s is stopping...", t.path)
			return
			//往kafka传信息
		case line, ok := <-t.tObj.Lines:
			if !ok {
				logrus.Warnf("tail file close reopen, filename:%s\n", t.path)
				time.Sleep(time.Second)
				continue
			}
			if len(strings.Trim(line.Text, "\r")) == 0 {
				continue
			}
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic
			msg.Value = sarama.StringEncoder(line.Text)
			kafkaProducer.ToMsgChan(msg)
		}
	}
}
