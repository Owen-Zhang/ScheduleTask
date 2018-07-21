package model

type WorkerFileConfig struct {
	Version  int    `json:"version"`
	FileName string `json:"filename"`
}

//文件服务器相关配制信息
type FileServerInfo struct {
	Hosts  string 
	Port   int    
}

//客户端相关信息
type WorkerInfo struct {
	Identification int   //标识，这个标识由服务器端统一分配，标识是哪个客户端
}