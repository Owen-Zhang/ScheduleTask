package ctrl

import (
	"ScheduleTask/taskworker/jobs"
	"ScheduleTask/model"
	"ScheduleTask/utils/system"
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"io"
	"encoding/json"
	"path"
	"sync"
)

const DataFolder = "Data/%s"

//运行任务(包括新增和重新启动)
func (this *Controller) start(request *Action) {
	//1: 查询数据，得到相关的实体数据
	task, err := this.Storage.GetTaskById(request.Id)
	if err != nil {
		fmt.Printf("start task has wrong: %s", err)
		return
	}

	if task == nil || task.Status != 1 {
		return
	}

	command := task.Command

	if task.TaskType == 1 {
		//生成文件夹(根文件夹,里面放有config文件和run文件夹)，run 放上传文件的解压内容
		taskfolder := strings.TrimSpace(task.RunFilefolder)
		datapath := fmt.Sprintf(DataFolder, taskfolder)

		//1: 检查上传文件是否有更新,如果没有更新就不用下载,没有配制文件也表示有更新
		configfile := fmt.Sprintf(`%s/config.txt`, datapath)
		if system.FileExist(configfile) {
			bytes, err := ioutil.ReadFile(configfile)
			if err != nil {
				fmt.Printf("read config file err: %s", err.Error())
				return
			}
			config := model.WorkerFileConfig{}
			errconfig := json.Unmarshal(bytes, &config)

			//反序列化有错就表示要重新下载文件及更新相应的配制信息
			if errconfig != nil {
				if err := this.updateConfig(task); err != nil {
					fmt.Printf("updateConfig err: %s", err.Error())
					return
				}
				
				if err := this.updateFileInfo(task); err != nil {
					fmt.Printf("updateFileInfo err: %s", err.Error())
					return
				}
			}

			//客户端版本比服务器低并且上传的zip文件名不同时，需要更新文件以及更新配制
			zipFileName := path.Base(task.ZipFilePath)
			if config.Version < task.Version && config.FileName != zipFileName {
				if err := this.updateConfig(task); err != nil {
					fmt.Printf("updateConfig err: %s", err.Error())
					return
				}
				
				if err := this.updateFileInfo(task); err != nil {
					fmt.Printf("updateFileInfo err: %s", err.Error())
					return
				}
			}

		} else {
			//新建文件夹，下载文件，新增配制信息
			if err := this.updateConfig(task); err != nil {
				fmt.Printf("updateConfig err: %s", err.Error())
					return
				}
				
				if err := this.updateFileInfo(task); err != nil {
					fmt.Printf("updateFileInfo err: %s", err.Error())
					return
				}
		}
		command = fmt.Sprintf("%s/%s/Run/%s", system.GetCurrentPath(), datapath, task.Command)
	}

	if jobs.ExistJob(task.Id) {
		jobs.RemoveJob(task.Id)
	}
	
	jobs.AddJob(&model.Task{
		Id				: task.Id,
		Name			: task.Name,
		CronSpec		: task.CronSpec,
		Command			: command,
		TaskType		: task.TaskType,
		TaskApiMethod 	: task.TaskApiMethod,
		TaskApiUrl 		: task.TaskApiUrl,
		Concurrent 		: task.Concurrent,
		TimeOut 		: task.TimeOut,
	})
}

//馬上運行任務
func (this *Controller) run(id int) {
	jobs.RunJob(id)
}

//停止任务
func (this *Controller) stop(id int) {
	jobs.RemoveJob(id)
}

//删除任务
func (this *Controller) delete(id int) {
	if jobs.ExistJob(id) {
		jobs.RemoveJob(id)
	}
	//最好還要刪除文件夾相關的東西
}

//更新配制文件
func (this *Controller) updateConfig(task * model.TaskExend) error {
	datapath := fmt.Sprintf(DataFolder, task.RunFilefolder)
	if !system.FileExist(datapath) {
			//数据文件夹没有，需要创建相关的文件夹
			if err := os.MkdirAll(datapath, 0777); err != nil {
				fmt.Printf("create run fileFolder err : %s", err.Error())
				return err
			}
	}
	
	configbytes, errjson := 
	json.Marshal(&model.WorkerFileConfig {
		Version 	: task.Version,
		FileName 	: path.Base(task.ZipFilePath),
	})
	
	if errjson != nil {
		return errjson
	}
	
	configfile := fmt.Sprintf("%s/config.txt", datapath)
	m := new(sync.Mutex)
	m.Lock()
	f, err := os.OpenFile(configfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(string(configbytes))
	defer m.Unlock()
	
	return nil
}

//下载文件，更新文件夹中的内容
func (this *Controller) updateFileInfo(task *model.TaskExend) error {
	datapath := fmt.Sprintf(DataFolder, task.RunFilefolder)
	tempzipfilefolder := fmt.Sprintf("%s/TempFile", datapath)
	runfilefolder := fmt.Sprintf("%s/Run", datapath)
	
	if !system.FileExist(tempzipfilefolder) {
		if err := os.MkdirAll(tempzipfilefolder, 0777); err != nil {
			fmt.Printf("create run TempFile err : %s", err.Error())
			return err
		}
	}
	
	if !system.FileExist(runfilefolder) {
		if err := os.MkdirAll(runfilefolder, 0777); err != nil {
			fmt.Printf("create run filefolder err : %s", err.Error())
			return err
		}
	}
	
	//下载文件
	fileserveraddrss := fmt.Sprintf("http://%s:%d/%s", this.FileServer.Hosts, this.FileServer.Port, task.ZipFilePath)
	res, errget := http.Get(fileserveraddrss)
	if errget != nil {
		fmt.Printf("DownLoad File err: %s\n",errget.Error())
		return errget
	}

	zipfile := fmt.Sprintf("%s/%s", tempzipfilefolder, path.Base(task.ZipFilePath))
	file, errzip := os.Create(zipfile)
	if errzip != nil {
		fmt.Printf("save temp file err: %s\n", errzip.Error())
		return errzip
	}

	defer file.Close()
	if _,err := io.Copy(file, res.Body); err != nil {
		fmt.Printf("copy temp file err: %s\n", err.Error())
		return err
	}
	defer res.Body.Close()

	//解压到指定的文件夹中
	if err := system.UnzipFile(zipfile, runfilefolder); err != nil {
		fmt.Printf("unzipfile has wrong err: %s", err.Error())
		return err
	}

	defer os.Remove(zipfile)
	
	return nil
}

//加载本地的任务到任务队列中
func (this *Controller) AddAutoRunSelfTask(identification int) {
	list,_ := this.Storage.TaskGetList(1, 100, 1, 0, identification)
	if list == nil {
		return
	}
	
	for _, val := range list {
		this.Actionlist <- Action {
			ActionType : 2,
			Id         : val.Id,
		}
	}
}