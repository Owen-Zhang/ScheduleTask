package main

import (
	"os"
	"ScheduleTask/taskworker/server"
	"ScheduleTask/utils/system"
	"github.com/astaxie/beego/logs"
	"ScheduleTask/taskworker/health"
	_"ScheduleTask/taskworker/etc"
)

//可以改为客户端报告，带上相应的客户端信息,让服务端去检查客户端的机子状态,以及分配相应的任务

func main() {

	logs.SetLogger(logs.AdapterFile,`{"filename":"worker.log"}`)

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
		}
	}()

	work, err := server.NewWorker()
	if err != nil {
		panic(err)
		os.Exit(-1)
	}

	defer func() {
		work.Stop()
		os.Exit(0)
	}()

	go func() {
		health.Heartbeat()
	}()

	if err := work.Start(); err != nil {
		panic(err)
		os.Exit(-2)
	}

	system.InitSignal(nil)
}
