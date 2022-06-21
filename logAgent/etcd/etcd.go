package etcd

import (
	"ch06-qimiProject/logAgent/common"
	"ch06-qimiProject/logAgent/tailfile"
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"time"
)

var client *clientv3.Client

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		logrus.Error("connent etcd failed:", err)
		return
	}
	logrus.Info("init etcd success")
	return
}

func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Error("getConf failed:", err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Error("长度为0，没有数据")
		return
	}
	ret := resp.Kvs[0]
	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func WatchConf(key string) {
	for {
		watchCh := client.Watch(context.Background(), key)
		for wresp := range watchCh {
			logrus.Info("get new conf from etcd")
			for _, evt := range wresp.Events {
				var newConf []common.CollectEntry
				//如果是删除
				if evt.Type == clientv3.EventTypeDelete {
					tailfile.SendNewConf(newConf)
					continue
				}
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Error(err)
					continue
				}
				tailfile.SendNewConf(newConf)
			}
		}
	}
}
