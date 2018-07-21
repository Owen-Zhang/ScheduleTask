package model

//总体返回的数据结构
type WorkerResponse struct {
	Success bool        `json:"success"` //是否成功
	Message string      `json:"message"` //消息内容
	Data    interface{} `json:"data"`    //业务数据
}
