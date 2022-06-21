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

	//put
	str := `[{"path":"D:/goprojects/src/ch06-qimiProject/demo/tailf-demo/xx.log","topic":"abc"}]`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "collect_log_172.22.106.111_conf", str)
	if err != nil {
		fmt.Println(err)
		return
	}
	cancel()

	//get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	gr, err := cli.Get(ctx, "collect_log_172.22.106.111_conf")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(gr.Kvs[0])
	for _, ev := range gr.Kvs {
		fmt.Println(string(ev.Key))
		fmt.Println(string(ev.Value))
	}

	cancel()
}
