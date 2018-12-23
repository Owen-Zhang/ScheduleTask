package controllers

import (
	"ScheduleTask/model"
	"strings"
	"ScheduleTask/utils/system"
	"github.com/astaxie/beego/logs"
	"ScheduleTask/taskserver/app/healthy"
)

type WorkerController struct {
	BaseController
}


func (this *WorkerController) List() {
	list, _ := dataaccess.GetWorkerList(2)

	//将客户端报告信息填充到列表中
	for _,value := range  healthy.Health.WorkerList  {
		for _, temp := range list {
			if value.WorkerInfo.WorkerKey == temp.Key {
				temp.Status = 3
				break;
			}
		}
	}

	this.Data["pageTitle"] = "worker列表"
	this.Data["list"] = list
	this.display()
}


// 新增
func (this *WorkerController) Add() {
	this.Data["pageTitle"] = "添加worker"
	this.display()
}

func (this *WorkerController) SaveWork() {
	worker := new(model.Worker)

	worker.Name = strings.TrimSpace(this.GetString("worker_name"))
	worker.Note = strings.TrimSpace(this.GetString("worker_note"))
	logs.Info(worker.Note)
	worker.Key = system.GetUuid();
	worker.Status = 0

	temp, _ := dataaccess.GetOneWorker(worker.Name, "", 0)
	if temp != nil {
		this.ajaxMsg("已存在相同的worker名称", MSG_ERR)
	}

	err := dataaccess.NewWorker(worker)
	if err != nil {
		this.ajaxMsg(err.Error(), MSG_ERR)
	}
	this.ajaxMsg("", MSG_OK)
}

// 查看
func (this *WorkerController) View() {
	id, _ := this.GetInt("id")

	worker, err := dataaccess.GetOneWorker("", "", id)
	if err != nil {
		this.showMsg(err.Error())
	}

	this.Data["pageTitle"] = "查看worker"
	this.Data["worker"] = worker
	this.Data["isview"] = 1
	this.display("worker/add")
}