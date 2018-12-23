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

	var info system.WorkerInfo
	if err := json.Unmarshal(by, &info); err != nil {
		logs.Error("worker and center heate has wrong: %s", err.Error())
		this.ajaxMsg(err.Error(), MSG_ERR)
	}

	if (info.WorkerKey == "") {
		this.ajaxMsg("请传入worker的标识(key)", MSG_ERR)
	}

	//判断是否是生成的worker信息，key
	worker, err := dataaccess.GetOneWorker("", info.WorkerKey, 0)
	if (err != nil) {
		logs.Error("客户报告状态检查失败: %s", err.Error())
		this.ajaxMsg("中心检查出现错误，请联系管理人员", MSG_ERR)
	}

	if (worker == nil) {
		this.ajaxMsg("传入worker的标识(key)中心不能识别", MSG_ERR)
	}

	healthy.Input(&info)
	this.ajaxMsg("ok", MSG_OK)
}