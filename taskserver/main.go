package main

import (
	"net/http"
	"html/template"
	"github.com/astaxie/beego"
	"ScheduleTask/taskserver/app/controllers"
	"ScheduleTask/taskserver/app/healthy"
	"ScheduleTask/storage"
	"fmt"
)
//还要增加一个任务状态，表示客户机出了问题，不是人为的结束任务的运行
//worker机器要存一个zip文件名， 用来判断本地和服务器上的文件是否相同，不同的话就要更新
func main()  {

	defer func(){
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	arg := &storage.DataStorageArgs{
		Hosts: beego.AppConfig.String("db.host"),
		DBName: beego.AppConfig.String("db.name"),
		User: beego.AppConfig.String("db.user"),
		Password: beego.AppConfig.String("db.password"),
		Port: beego.AppConfig.DefaultInt("db.port", 3306),
	}
	dataaccess, errData := storage.NewDataStorage(arg)
	if errData != nil {
		fmt.Printf("init storage dataaccess has wrong: %s", errData)
		return
	}
	defer dataaccess.Close()

	//1: 加载要执行的任务数据(多久去check worker的状态)
	healthy.InitHealthCheck(beego.AppConfig.String("site.cron"), dataaccess)

	// 设置默认404页面
	beego.ErrorHandler("404", func(rw http.ResponseWriter, r *http.Request) {
		t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/404.html")
		data := make(map[string]interface{})
		data["content"] = "page not found"
		t.Execute(rw, data)
	})

	controllers.InitCtrl(dataaccess)
	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/profile", &controllers.MainController{}, "*:Profile")
	beego.Router("/gettime", &controllers.MainController{}, "*:GetTime")
	beego.Router("/help", &controllers.HelpController{}, "*:Index")
	beego.AutoRouter(&controllers.TaskController{})
	beego.AutoRouter(&controllers.GroupController{})
	beego.AutoRouter(&controllers.WorkerController{})

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}
