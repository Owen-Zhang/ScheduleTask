package system

import (
	"os"
	"runtime"
	"github.com/astaxie/beego/logs"
	"time"
)

//心跳检查的客户端信息(worker)
type HealthInfo struct {
	WorkerInfo    WorkerInfo   `json:"workerinfo"`                //worker机器信息
	Status     	  int                              			   	  // 是否可用
	TimeOut       time.Time                                       //发送时间
}

// WorkerInfo 客户机器信息
type WorkerInfo struct {
	WorkerKey              string        `json:"workerkey"`            //服务端生成的worker标识
	Name       			   string        `json:"name"`                 // 机器名称
	Ip         			   string        `json:"ip"`                   // 地址
	Port                   string        `json:"port"`                 // 端口号
	OsName 			       string        `json:"osname"`               // 系统信息(windows, linux,...)
	Note       			   string        `json:"note"`                 // 说明
	Memory                 MemoryInfo                       		   // 内存相关信息
}

type CpuInfo struct {

}

type MemoryInfo struct {
	Total int
	Free  int
	UsedPercent float64
}

type DiskInfo struct {

}

var SystemInfo *WorkerInfo

func init()  {
	hostName, _ := os.Hostname()
	ip, err := GetIntranetIp()
	if (err != nil) {
		logs.Error(err)
	}

	SystemInfo = &WorkerInfo{
		Name: hostName,
		Ip  : ip,
		Port: "8985",
		OsName: runtime.GOOS,
	}
}
