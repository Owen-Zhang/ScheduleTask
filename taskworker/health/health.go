package health

import (
	"time"
	"fmt"
	"strings"
	"github.com/imroc/req"
	"ScheduleTask/utils/system"
	"ScheduleTask/taskworker/etc"
	"github.com/astaxie/beego/logs"
)

// Heartbeat 定时向主机报告当前机器的一些运行情况
func Heartbeat() {
	var duration time.Duration = 10 * time.Second
	ticker := time.NewTicker(duration)

	for {
		select {
		case <-ticker.C:
			workerApiConfig := etc.GetApiServerArg()

			worker := system.SystemInfo
			worker.WorkerKey = workerApiConfig.Key
			worker.Ip = workerApiConfig.Ip
			worker.Port = string([]rune(strings.TrimSpace(workerApiConfig.Bind))[1:])
			worker.Note = workerApiConfig.Note

			centerInfo := etc.GetCenterInfo()
			req.SetTimeout(10 * time.Second)
			_, err := req.Post(
							fmt.Sprintf("http://%s:%d/health/ping", centerInfo.Hosts, centerInfo.Port),
							req.BodyJSON(worker))
			logs.Error(err)
		}
	}
}