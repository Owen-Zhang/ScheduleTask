package model

type Fileinfo struct {
	FilePath 		string `json:"filepath"`       //文件路径，如：order/detail
	FileSuffixName  string `json:"filesuffixname"` //文件的后缀名 如: exe, jpg
	FileContent  	string `json:"filecontent"`    //文件的byte的64编码
}

type FileResponse struct {
	Status      		bool     `json:"status"` 		  //状态，如：true,false
	Message     		string   `json:"message"` 		  //信息，如：提示信息
	FileName     		string   `json:"filename"` 		  //文件名称，如：4545489.exe
}

const WorkerUrl = "http://%s:%d/worker/%s"