package model

//心跳检查的客户端信息(worker)
type HealthInfo struct {
	Id 		   int                              // 任务ID
	Name       string                           // 机器名称
	Url        string                           // 地址
	Port       int                              // 端口
	SystemInfo string                           // 系统信息(windows, linux,...)
	Note       string                           // 说明
	Status     int                              // 是否可用
}
