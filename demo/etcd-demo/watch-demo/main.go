package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cli.Close()

	watchCh := cli.Watch(context.Background(), "collect_log_172.22.106.111_conf")
	fmt.Println(watchCh)
	for wresp := range watchCh {
		fmt.Println(watchCh)
		for _, evt := range wresp.Events {
			fmt.Println(evt)
		}
	}
}
