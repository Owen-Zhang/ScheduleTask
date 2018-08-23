package model

import "runtime"

const (
	ServerTempFileFolder  = "temp"
	WorkerRunDir   = "data" //所有的文件任务文件放在此目录下
	WorkerFileRunDir = "run" //将用户上传的文件解压后放的目录
	CNTimeFormat = "2006-01-02 15:04:05"
	SystemWindows = "windows"
	SystemLinux = "linux"
	WindowsShellExt = "bat"
	LinuxShellExt = "sh"
)

var Common *CommonInfo

type WorkerFileConfig struct {
	Version  int    `json:"version"`
	FileName string `json:"filename"`
}

//文件服务器相关配制信息
type FileServerInfo struct {
	Hosts  string 
	Port   int    
}

type CommonInfo struct {
	SystemName string   //系统名称(windows, linux)
}

func init() {
	Common = &CommonInfo{
		SystemName : runtime.GOOS,
	}
}