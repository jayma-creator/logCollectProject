package tailfile

import (
	"ch06-qimiProject/logAgent/common"
	"fmt"
	"github.com/sirupsen/logrus"
)

type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask
	collectEntryList []common.CollectEntry
	confChan         chan []common.CollectEntry
}

var ttMgr *tailTaskMgr

func Init(allConf []common.CollectEntry) (err error) {
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20), //存放已创建的tailtask
		collectEntryList: allConf,
		confChan:         make(chan []common.CollectEntry), //初始化新配置的管道
	}
	for _, conf := range allConf {
		//对每一个配置项创建一个日志收集任务
		tt := newTailTask(conf.Path, conf.Topic)
		err = tt.Init()
		if err != nil {
			fmt.Println(err)
			return
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		//每新建一个tailtask任务都存放在map里
		ttMgr.tailTaskMap[tt.path] = tt
		//往管道发信息到kafka
		go tt.run()
	}
	//监控有无新配置
	go ttMgr.watch()
	logrus.Info("init tailfile success")
	return
}

func (t *tailTaskMgr) watch() {
	for {
		//读取新配置
		newConf := <-ttMgr.confChan
		logrus.Infof("get new conf from etcd, conf:%v", newConf)

		for _, conf := range newConf {
			//如果新配置包含旧的，不作处理
			if t.isExist(conf) {
				continue
			}
			//如果新配置里有新的，旧配置没有，新建一个tailtask任务
			tt := newTailTask(conf.Path, conf.Topic)
			err := tt.Init()
			if err != nil {
				fmt.Println(err)
				return
			}
			logrus.Infof("create a tail task for path:%s success", conf.Path)
			go tt.run()
		}
		//如果新配置里没有，旧的有，把旧的停掉
		for key, task := range t.tailTaskMap {
			var found bool
			for _, conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				logrus.Infof("the task collect path:%s need to stop.", task.path)
				delete(t.tailTaskMap, key)
				task.cancel()
			}
		}
	}
}

func (t *tailTaskMgr) isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}

//发送配置项
func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
