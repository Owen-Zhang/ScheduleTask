package healthy

/*
import (
	"github.com/Owen-Zhang/cron"
	"sync"
	"ScheduleTask/model"
	"ScheduleTask/storage"
)

var (
	mainCron 	*cron.Cron
	lock     	sync.Mutex
	dataaccess  *storage.DataStorage
)

func init()  {
	mainCron = cron.New()
	mainCron.Start()
}


//加载worker，监控worker
func InitHealthCheck(spec string, access *storage.DataStorage) {
	dataaccess = access
	list, err := dataaccess.GetWorkerList(1, "")
	if err != nil || list == nil || len(list) == 0 {
		return
	}

	for _, heal := range list {
		AddHealthyCheck(spec, heal)
	}
}

//增加心跳任务，检查worker机子是否正常运行
func AddHealthyCheck(spec string, info *model.HealthInfo) bool {
	heal, err := newHealth(info)
	if err != nil {
		return false
	}

	lock := sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	if ExistJob(heal.id) {
		return false
	}

	if err := mainCron.AddJob(spec, heal); err == nil {
		return true
	}
	return false
}

//删除运行中的任务
func RemoveJob(id int) {
	if !ExistJob(id) {
		return
	}
	mainCron.RemoveJob(func(e *cron.Entry) bool {
		if v, flag := e.Job.(*Health); flag {
			if v.id == id {
				return true
			}
		}
		return false
	})
}

//判断任务是否在指行队列中
func ExistJob(id int) bool  {
	entries := mainCron.Entries()
	for _, e := range entries {
		if v, flag := e.Job.(*Health); flag {
			if v.id == id {
				return true
			}
		}
	}
	return false
}
*/