package api

import (
	"net/http"
	"ScheduleTask/taskworker/ctrl"
	"ScheduleTask/model"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

//加入任务
func (this *ApiServer) starttask(c *gin.Context) {
	response := &model.WorkerResponse{
		Success: false,
		Message: "",
	}

	idtemp := c.PostForm("id")
	id, err :=strconv.Atoi(idtemp)
	if err != nil || id <= 0 {
		response.Message = "please input right task id"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	this.controller.Actionlist <- ctrl.Action{
		ActionType: 2,
		Id:         id,
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}

// 批量加入任务
func (this *ApiServer) batchstarttask(c *gin.Context) {
	response := &model.WorkerResponse{
		Success: false,
		Message: "",
	}

	ids := c.PostForm("ids")
	if ids == "" {
		response.Message = "please input task ids"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	idArr := strings.Split(strings.TrimSpace(ids), ",")
	for _,val := range idArr {
		id, err :=strconv.Atoi(val)
		if err != nil || id <= 0 {
			continue
		}
		this.controller.Actionlist <- ctrl.Action{
			ActionType : 2,
			Id         : id,
		}
	}
	response.Success = true
	c.JSON(http.StatusOK, response)
}

//運行任務
func (this *ApiServer) runtask(c *gin.Context) {
	response := &model.WorkerResponse{
		Success: false,
		Message: "",
	}

	idtemp := c.PostForm("id")
	id, err :=strconv.Atoi(idtemp)
	if err != nil || id <= 0 {
		response.Message = "please input right task id"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	this.controller.Actionlist <- ctrl.Action{
		ActionType: 5,
		Id:         id,
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}

//停止任务
func (this *ApiServer) stoptask(c *gin.Context) {
	response := &model.WorkerResponse{
		Success: false,
		Message: "",
	}

	idtemp := c.PostForm("id")
	id, err :=strconv.Atoi(idtemp)
	if err != nil || id <= 0 {
		response.Message = "please input right task id"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	this.controller.Actionlist <- ctrl.Action{
		ActionType: 3,
		Id:         id,
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}

//删除任务
func (this *ApiServer) deletetask(c *gin.Context) {
	response := &model.WorkerResponse{
		Success: false,
		Message: "",
	}

	idtemp := c.PostForm("id")
	id, err :=strconv.Atoi(idtemp)
	if err != nil || id <= 0 {
		response.Message = "please input right task id"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	this.controller.Actionlist <- ctrl.Action{
		ActionType: 4,
		Id:         id,
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}
