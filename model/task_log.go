package model

//任务执行日志
type TaskLog struct {
	Id 			int				  //日志主键
	TaskId  	int				  //任务主键
	Output      string            //正常输出值
	Error		string			  //错误输出值
	ProcessTime int				  //执行时间
	CreateTime  int64			  //创建时间
	Status      int				  //日志状态
}
