package main

import (
	"fmt"
	"net/http"
	"html/template"
	"ScheduleTask/storage"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"ScheduleTask/taskserver/app/controllers"
	"ScheduleTask/taskserver/app/healthy"
)


func main()  {

	logs.SetLogger(logs.AdapterFile,`{"filename":"server.log"}`)

	defer func(){
		if err := recover(); err != nil {
			logs.Error(err)
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
		logs.Error(fmt.Sprintf("init storage dataaccess has wrong: %s", errData))
		return
	}
	defer dataaccess.Close()

	// 监控worker的报告状态
	go healthy.CheckWorkerStatus(dataaccess)

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
	beego.AutoRouter(&controllers.HealthController{})

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}
