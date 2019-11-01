package jobs

import (
	"ScheduleTask/model"
	"ScheduleTask/storage"
	"sync"

	"github.com/Owen-Zhang/cron2"
)

type CronArg struct {
	PoolSize int32
}

var (
	workPool chan bool
	mainCron *cron.Cron
	data     *storage.DataStorage
)

func NewCron(arg *CronArg, storage *storage.DataStorage) {
	if arg.PoolSize > 0 {
		workPool = make(chan bool, arg.PoolSize)
	}

	data = storage
	mainCron = cron.New()
	mainCron.Start()
}

//增加任务
func AddJob(task *model.Task) bool {
	job, err := newJobFromTask(task)
	if err != nil {
		return false
	}

	lock := sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	if ExistJob(job.id) {
		return false
	}

	if err := mainCron.AddJob(task.CronSpec, job); err == nil {
		return true
	}
	return false
}

//RemoveJob 删除运行中的任务
func RemoveJob(id int) {
	if !ExistJob(id) {
		return
	}
	mainCron.RemoveJob(func(e *cron.Entry) bool {
		if v, flag := e.Job.(*Job); flag {
			if v.id == id {
				return true
			}
		}
		return false
	})
}

//判断任务是否在指行队列中
func ExistJob(id int) bool {
	entries := mainCron.Entries()
	for _, e := range entries {
		if v, flag := e.Job.(*Job); flag {
			if v.id == id {
				return true
			}
		}
	}
	return false
}

// 運行任務中的job
func RunJob(id int) {
	entry := getEntryById(id)
	if entry != nil {
		entry.Job.Run()
	}
}

// 查詢出cron的
func getEntryById(id int) *cron.Entry {
	entries := mainCron.Entries()
	for _, e := range entries {
		if v, flag := e.Job.(*Job); flag {
			if v.id == id {
				return e
			}
		}
	}
	return nil
}
