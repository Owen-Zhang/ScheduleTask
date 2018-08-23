package controllers

import (
	"io/ioutil"
	"encoding/json"
	"ScheduleTask/utils/system"
	"github.com/astaxie/beego/logs"
	"ScheduleTask/taskserver/app/healthy"
)

type HealthController struct {
	BaseController
}

func (this * HealthController) Ping() {
	by, errbyte := ioutil.ReadAll(this.Ctx.Request.Body)
	if errbyte != nil {
		this.ajaxMsg(errbyte.Error(), MSG_ERR)
	}

	var info system.HealthInfo
	if err := json.Unmarshal(by, &info); err != nil {
		logs.Error("worker and center heate has wrong: %s", err.Error())
		this.ajaxMsg(err.Error(), MSG_ERR)
	}

	healthy.Input(&info)
	this.ajaxMsg("ok", MSG_OK)
}