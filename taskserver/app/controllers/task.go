package controllers

import (
	"time"
	"os"
	"strings"
	"strconv"
	"ScheduleTask/model"
	"ScheduleTask/utils/system"
	"fmt"
	"net/http"

	"ScheduleTask/taskserver/app/libs"
	"github.com/astaxie/beego"
	"ScheduleTask/taskserver/app/models/response"
	"github.com/Owen-Zhang/cron"
	"github.com/imroc/req"
	"io/ioutil"
	"encoding/base64"
	"errors"
)

type TaskController struct {
	BaseController
}

const tempFileFolder  = "TempFile"
//const workerUrl = "http://%s:%d/worker/%s"

// 任务列表
func (this *TaskController) List() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	groupId, _ := this.GetInt("groupid")
	workerid, _ := this.GetInt("workerid")

	// 分组列表
	groups, _ := dataaccess.TaskGroupGetList(1, 100)
	workersTemp, _:= dataaccess.GetWorkerList(2, "")
	
	result, count := dataaccess.TaskGetList(page, this.pageSize, -1, groupId, workerid)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.Name
		row["cron_spec"] = v.CronSpec
		row["status"] = v.Status
		row["description"] = v.Description
		row["tasktype"] = v.TaskType
		row["groupname"] = ""
		row["workname"] = ""
		
		for _,val := range groups {
			if val.Id == v.GroupId {
				row["groupname"] = val.GroupName
			}
		}	
		for _,val := range workersTemp {
			if val.Id == v.WorkerId {
				row["workname"] = val.Name
			}
		}
			
		list[k] = row
	}

	workers, _:= dataaccess.GetWorkerList(1, "")
	
	this.Data["pageTitle"] = "任务列表"
	this.Data["list"] = list
	this.Data["groups"] = groups
	this.Data["workers"] = workers
	this.Data["groupid"] = groupId
	this.Data["workerid"] = workerid
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("TaskController.List", "groupid", groupId, "workerid", workerid), true).ToString()
	this.display()
}

func (this *TaskController) UploadRunFile() {
	f, h, err := this.GetFile("files[]")
	defer f.Close()

	uploadResult := &response.ResultData{
		IsSuccess: false,
	}

	if err != nil {
		uploadResult.Msg = "请选择要上传的文件"
		this.jsonResult(uploadResult)

	} else {
		fileTool := &libs.FileTool{Url: h.Filename}
		exts := []string{"zip"}
		if !fileTool.CheckFileExt(exts) {
			uploadResult.Msg = "请上传正确的文件类型"
			this.jsonResult(uploadResult)
		}

		uuidFileName := fileTool.CreateUuidFile()
		if uuidFileName == "" {
			uploadResult.Msg = "文件保存出错，请重新选择文件"
			this.jsonResult(uploadResult)
		}

		filePath := tempFileFolder + "/" + uuidFileName
		os.Mkdir(tempFileFolder, 0777)

		if err := this.SaveToFile("files[]", filePath); err != nil {
			uploadResult.Msg = err.Error()
			this.jsonResult(uploadResult)
		}

		uploadResult.IsSuccess = true
		uploadResult.Data = &response.UploadFileInfo{
			OldFileName: fileTool.Url,
			NewFileName: uuidFileName,
		}
		this.jsonResult(uploadResult)
	}
}

// 添加任务
func (this *TaskController) Add() {
	groups, _ := dataaccess.TaskGroupGetList(1, 100)
	workers,_ := dataaccess.GetWorkerList(1, "")

	this.Data["groups"] = groups
	this.Data["workers"] = workers
	this.Data["pageTitle"] = "添加任务"
	this.display()
}

// 编辑任务
func (this *TaskController) Edit() {
	id, _ := this.GetInt("id")

	task, err := dataaccess.GetTaskById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	// 分组列表
	groups, _ := dataaccess.TaskGroupGetList(1, 100)
	workers,_ := dataaccess.GetWorkerList(1, "")

	this.Data["groups"] = groups
	this.Data["workers"] = workers
	this.Data["task"] = task
	this.Data["pageTitle"] = "编辑任务"
	
	status := this.GetString("status", "")
	this.Data["status"] = status

	this.display("task/add")
}

//保存任务
func (this *TaskController) SaveTask() {
	id, _ := this.GetInt("id", 0)
	isNew := true
	if id != 0 {
		isNew = false
	}

	task := new(model.TaskExend)
	if !isNew {
		var err error
		task, err = dataaccess.GetTaskById(id)
		if err != nil {
			this.showMsg(err.Error())
		}
	} else {
		task.UserId = this.userId
	}

	if isNew {
		task.TaskType, _ = this.GetInt("task_type")
	}
		
	task.Name = strings.TrimSpace(this.GetString("task_name"))
	task.Description = strings.TrimSpace(this.GetString("description"))
	task.GroupId, _ = this.GetInt("group_id")
	task.Concurrent, _ = this.GetInt("concurrent")
	task.CronSpec = strings.TrimSpace(this.GetString("cron_spec"))
	task.Command = strings.TrimSpace(this.GetString("command"))
	task.Notify, _ = this.GetInt("notify")
	task.TimeOut, _ = this.GetInt("timeout")
	task.WorkerId,_ =  this.GetInt("worker_id")
	task.TaskApiMethod = ""
	if task.TaskType == 2 {
		task.TaskApiMethod = strings.TrimSpace(this.GetString("task_method"))
	}
	task.TaskApiUrl = strings.TrimSpace(this.GetString("task_url"))
	task.ApiHeader = strings.TrimSpace(this.GetString("api_header"))
	useruploadfile := strings.TrimSpace(this.GetString("oldzipfile"))
	
	isUploadNewFile := false
	if task.TaskType == 1 && task.OldZipFile != useruploadfile {
		isUploadNewFile = true
		task.OldZipFile = useruploadfile
	}
	notifyEmail := strings.TrimSpace(this.GetString("notify_email"))

	resultData := &response.ResultData{IsSuccess: false, Msg: ""}
	if notifyEmail != "" {
		emailList := make([]string, 0)
		tmp := strings.Split(notifyEmail, ";")
		for _, v := range tmp {
			v = strings.TrimSpace(v)
			if !libs.IsEmail([]byte(v)) {
				resultData.Msg = "无效的Email地址：" + v
				this.jsonResult(resultData)
			} else {
				emailList = append(emailList, v)
			}
		}
		task.NotifyEmail = strings.Join(emailList, ";")
	}

	if task.Name == "" || task.CronSpec == "" || 
		((task.TaskType == 0 || task.TaskType == 1) && task.Command == "") || 
		(task.TaskType == 2 && (task.TaskApiUrl == "" || task.TaskApiMethod == ""))  {
		resultData.Msg = "请填写完整信息"
		this.jsonResult(resultData)
	}
	
	if _, err := cron.Parse(task.CronSpec); err != nil {
		resultData.Msg = "cron表达式无效"
		this.jsonResult(resultData)
	}

	if task.TaskType == 1 && isUploadNewFile && task.OldZipFile != "" {
		
		/*runFileName: 记录处理过的文件名(为了保存文件名不重复，重新取文件名); OldZipFile: 用户上传的文件名*/
		runFileName := strings.TrimSpace(this.GetString("runfilename"))

		filepath := tempFileFolder + "/" +  runFileName
		//上传文件到文件服务器
		if system.IsExist(filepath) {
			filename, err := this.uploadfile(filepath)
			if err != nil {
				resultData.Msg = err.Error()
				this.jsonResult(resultData)
			}
			task.ZipFilePath = filename
		} else {
			resultData.Msg = fmt.Sprintf("TempFile/%s is not exists", runFileName)
			this.jsonResult(resultData)
		}
		if task.RunFilefolder == "" {
			task.RunFilefolder = system.GetUuid()
		}
	}

	//保存数据库
	if isNew {
		task.Version = 1
		if err := dataaccess.TaskAdd(task); err != nil {
			resultData.Msg = err.Error()
			this.jsonResult(resultData)
		}
	} else {
		task.Version += 1
		if err := dataaccess.UpdateFrontTask(task); err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}
	}

	resultData.IsSuccess = true
	this.jsonResult(resultData)
}

// uploadfile 上传文件到文件服务器
func (this *TaskController) uploadfile(filename string) (string, error) {
	url := fmt.Sprintf("http://%s:%s/upload",
		beego.AppConfig.String("file.host"),
		beego.AppConfig.String("file.port"))

	fileopen, err1 := os.Open(filename)
	if err1 != nil {
		fmt.Println(err1.Error())
		return "", err1
	}
	defer fileopen.Close()

	fd,err2 := ioutil.ReadAll(fileopen)
	if err2 != nil {
		fmt.Println(err2.Error())
		return "", err2
	}
	encodeString := base64.StdEncoding.EncodeToString(fd)

	fileresponse, err :=
		req.Post(url, req.BodyJSON(&model.Fileinfo{
			FilePath: "job",
			FileSuffixName: "zip",
			FileContent: encodeString,
		}))
	if err != nil {
		return "", err
	}

	var res = &model.FileResponse{}
	fileresponse.ToJSON(res)
	if !res.Status {
		return "", errors.New(res.Message)
	}
	return res.FileName, nil
}

// 任务执行日志列表
func (this *TaskController) Logs() {
	taskId, _ := this.GetInt("id")
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	task, err := dataaccess.GetTaskById(taskId)
	if err != nil {
		this.showMsg(err.Error())
	}

	result, count := dataaccess.TaskLogGetList(page, this.pageSize, 1, task.Id)

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["status"] = v.Status
		list[k] = row
	}

	this.Data["pageTitle"] = "任务执行日志"
	this.Data["list"] = list
	this.Data["task"] = task
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("TaskController.Logs", "id", taskId), true).ToString()
	this.display()
}

// 查看日志详情
func (this *TaskController) ViewLog() {
	id, _ := this.GetInt("id")

	taskLog, err := dataaccess.TaskLogGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	task, err := dataaccess.GetTaskById(taskLog.TaskId)
	if err != nil {
		this.showMsg(err.Error())
	}

	data := make(map[string]interface{})
	data["id"] = taskLog.Id
	data["output"] = taskLog.Output
	data["error"] = taskLog.Error
	data["start_time"] = beego.Date(time.Unix(taskLog.CreateTime, 0), "Y-m-d H:i:s")
	data["process_time"] = float64(taskLog.ProcessTime) / 1000
	data["ouput_size"] = libs.SizeFormat(float64(len(taskLog.Output)))
	data["status"] = taskLog.Status

	this.Data["task"] = task
	this.Data["data"] = data
	this.Data["pageTitle"] = "查看日志"
	this.display()
}

// 批量操作日志
func (this *TaskController) LogBatch() {
	action := this.GetString("action")
	ids := this.GetStrings("ids")
	if len(ids) < 1 {
		this.ajaxMsg("请选择要操作的项目", MSG_ERR)
	}
	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		switch action {
		case "delete":
			dataaccess.TaskLogDelById(id)
		}
	}

	this.ajaxMsg("", MSG_OK)
}

// 启动任务
func (this *TaskController) Start() {
	result := &response.ResultData{
		IsSuccess: false,
		Msg:       "",
	}
	id, _ := this.GetInt("id")
	
	if id <= 0 {
		result.Msg = "请操作正常的任务"
		this.jsonResult(result)
	}

	task, err := dataaccess.GetTaskById(id)
	if err != nil {
		result.Msg = err.Error()
		this.jsonResult(result)
	}
	if task != nil {
		worker, err := dataaccess.GetOneWorker("", task.WorkerId)
		if err != nil {
			result.Msg = err.Error()
			this.jsonResult(result)
		}
		
		if worker != nil {
			posturl := fmt.Sprintf(model.WorkerUrl, worker.Url, worker.Port, "starttask")
			fmt.Println(posturl)

			updateerr := dataaccess.TaskUpdateStatus(id, 1)
			if updateerr != nil {
				result.Msg = updateerr.Error()
				this.jsonResult(result)
			}
			
			res, err :=
			req.Post(posturl, req.Param{"id" : id})
			
			if err != nil || res.Response().StatusCode != http.StatusOK {
				if err == nil {
					result.Msg = "[Start]通知客戶端失敗"
				} else {
					result.Msg = err.Error()
				}

				//将状态改回去
				dataaccess.TaskUpdateStatus(id, 0)

				this.jsonResult(result)
			}
		}
	}

	this.jsonResult(&response.ResultData{
		IsSuccess: true,
		Msg:       "",
		Data: &response.JobInfo{
			Status: 1,
			Prev:   "-",
			Next:   "-",
		},
	})
}

// 直接运行任务
func (this *TaskController) Run()  {
	result := &response.ResultData{
		IsSuccess: false,
		Msg:       "",
	}

	id, _ := this.GetInt("id")
	if id <= 0 {
		result.Msg = "请操作正常的任务"
		this.jsonResult(result)
	}

	task, err := dataaccess.GetTaskById(id)
	if err != nil {
		result.Msg = err.Error()
		this.jsonResult(result)
	}

	//这里可以分成两个动作，先start再run，后台统一去作处理
	if task.Status != 1 {
		result.Msg = "请先开始任务再运行"
		this.jsonResult(result)
	}

	if task != nil {
		worker, err := dataaccess.GetOneWorker("", task.WorkerId)
		if err != nil {
			result.Msg = err.Error()
			this.jsonResult(result)
		}

		if worker != nil {
			posturl := fmt.Sprintf(model.WorkerUrl, worker.Url, worker.Port, "stoptask")
			fmt.Println(posturl)

			res, err := req.Post(posturl, req.Param{"id": id})
			if err != nil || res.Response().StatusCode != http.StatusOK {
				result.Msg = err.Error()
				if err == nil {
					result.Msg = "[Stop]通知客戶端失敗"
				}
				this.jsonResult(result)
			}
		}
	}
	
	this.jsonResult(&response.ResultData{
		IsSuccess: true,
		Msg:       "",
		Data: &response.JobInfo{
			Status: 0,
			Prev:   "-",
			Next:   "-",
		},
	})
}

// 暂停任务
func (this *TaskController) Pause() {
	result := &response.ResultData{
		IsSuccess: false,
		Msg:       "",
	}	
	
	id, _ := this.GetInt("id")
	if id <= 0 {
		result.Msg = "请操作正常的任务"
		this.jsonResult(result)
	}
	
	task, err := dataaccess.GetTaskById(id)
	if err != nil {
		result.Msg = err.Error()
		this.jsonResult(result)
	}
	
	if task != nil {
		worker, err := dataaccess.GetOneWorker("", task.WorkerId)
		if err != nil {
			result.Msg = err.Error()
			this.jsonResult(result)
		}
		
		if worker != nil {
			posturl := fmt.Sprintf(model.WorkerUrl, worker.Url, worker.Port, "stoptask")
			fmt.Println(posturl)
			
			res, err :=
			req.Post(posturl, req.Param{"id" : id})
			
			if err != nil || res.Response().StatusCode != http.StatusOK {
				result.Msg = err.Error()
				if err == nil {
					result.Msg = "[Stop]通知客戶端失敗"
				}
				this.jsonResult(result)
			}
			
			if err := dataaccess.TaskUpdateStatus(id, 0); err != nil {
				result.Msg = err.Error()
				this.jsonResult(result)
			}
		}
	}
	
	this.jsonResult(&response.ResultData{
		IsSuccess: true,
		Msg:       "",
		Data: &response.JobInfo{
			Status: 0,
			Prev:   "-",
			Next:   "-",
		},
	})
}

// 删除任务
func (this *TaskController) Delete() {
	result := &response.ResultData{
		IsSuccess: false,
		Msg:       "",
		Data:      true,
	}	
	id, _ := this.GetInt("id")
	if id <= 0 {
		result.Msg = "请操作正常的任务"
		this.jsonResult(result)
	}

	task, err := dataaccess.GetTaskById(id)
	if err != nil {
		result.Msg = err.Error()
		this.jsonResult(result)
	}
	
	if task != nil {
		//根据任务去找到worker相关的地址信息
		worker, err := dataaccess.GetOneWorker("", task.WorkerId)
		if err != nil {
			result.Msg = err.Error()
			this.jsonResult(result)
		}
		if worker != nil {
			posturl := fmt.Sprintf(model.WorkerUrl, worker.Url, worker.Port, "deletetask")
			fmt.Println(posturl)
			
			res, err :=
			req.Post(posturl, req.Param{"id" : id})
			
			if err != nil || res.Response().StatusCode != http.StatusOK {
				result.Msg = err.Error()
				if err == nil {
					result.Msg = "[Delete]通知客戶端失敗"
				}
				this.jsonResult(result)
			}
			
			if err := dataaccess.TaskDel(id); err != nil {
				result.Msg = err.Error()
				this.jsonResult(result)
			}
		}	
	}
		
	result.IsSuccess = true
	this.jsonResult(result)
}
