package model

const (
	TASK_SUCCESS = 0  // 任务执行成功
	TASK_ERROR   = -1 // 任务执行出错
	TASK_TIMEOUT = -2 // 任务执行超时
)

//任务
type Task struct {
	Id          	int           //任务的主键
	TaskType    	int           //任务类型 0命令型，1运行本地文件(上传文件),2调用外部接口
	Name        	string        //任务名称
	CronSpec   		string     	  //cron表达式
	RunFilefolder   string    	  //任务的文件夹（代码放的文件夹名）
	OldZipFile    	string        //原来的zip文件
	Command     	string        //任务的命令如Init.exe xxx
	TaskApiUrl      string        //API地址，如有端口号需要加上端口
	TaskApiMethod   string        //提交方式(POST, GET)
	ApiHeader       string        //提交的header
	TimeOut     	int           //任务执行的超时时间
	Concurrent  	int   		  //是否允许在再一次没有运行完成的情况运行下一次
	Notify      	int           //是否需要通知
	NotifyEmail 	string 		  //通知的邮件地址
	Version     	int  		  //程序的版本号
	ZipFilePath     string 		  //zip的存储位置(只有相对路径,不包含ip和端口)
	WorkerId        int           //客户端的引用编号
}

//前端在使用
type TaskExend struct {
	Task
	UserId       int     //用户编号
	GroupId      int     //组
	Description  string  //描述
	//RunFileName  string  //运行的文件
	Status       int     //状态
	ExecuteTimes int     //运行时间
	PrevTime     int64   //上次运行的开始时间
	CreateTime   int64   //创建时间
}

