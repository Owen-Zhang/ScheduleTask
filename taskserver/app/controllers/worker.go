package controllers

/*
import (
	"ScheduleTask/model"
	"strings"
)

type WorkerController struct {
	BaseController
}


func (this *WorkerController) List() {
	list, _ := dataaccess.GetWorkerList(2, "")

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
	worker := new(model.HealthInfo)

	worker.Name = strings.TrimSpace(this.GetString("worker_name"))
	worker.SystemInfo = strings.TrimSpace(this.GetString("worker_systeminfo"))
	worker.Url = strings.TrimSpace(this.GetString("worker_url"))
	worker.Port,_ = this.GetInt("worker_port", 0)
	worker.Note = strings.TrimSpace(this.GetString("worker_note"))
	worker.Status = 1

	_, errT := dataaccess.GetOneWorker(worker.Name, 0)
	if errT == nil {
		this.ajaxMsg("已存在相同的worker名称", MSG_ERR)
	}

	err := dataaccess.AddWorker(worker)
	if err != nil {
		this.ajaxMsg(err.Error(), MSG_ERR)
	}
	this.ajaxMsg("", MSG_OK)
}
*/