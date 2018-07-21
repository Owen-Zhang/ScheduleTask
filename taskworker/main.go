package main

import (
	"ScheduleTask/taskworker/server"
	"os"
	"ScheduleTask/utils/system"
	"fmt"
)

//可以改为客户端报告，带上相应的客户端信息,让服务端去检查客户端的机子状态,以及分配相应的任务

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
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

	if err := work.Start(); err != nil {
		panic(err)
		os.Exit(-2)
	}

	system.InitSignal(nil)
}
