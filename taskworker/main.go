package main

import (
	"os"
	"ScheduleTask/taskworker/server"
	"ScheduleTask/utils/system"
	"github.com/astaxie/beego/logs"
	"ScheduleTask/taskworker/health"
	_"ScheduleTask/taskworker/etc"
)

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
