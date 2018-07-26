package jobs

import (
	"bytes"
	"fmt"
	"ScheduleTask/model"
	"os/exec"
	"time"

	"github.com/axgle/mahonia"
	"github.com/imroc/req"
	"io/ioutil"
	"path/filepath"
	"os"
	"errors"
	"net/http"
	"strings"
)

type Job struct {
	id 		   int                                               // 任务ID
	name       string                                            // 任务名称
	task       *model.Task                                       // 任务对象
	runFunc    func(time.Duration) (string, string, error, bool) // 执行函数
	status     int                                               // 任务状态，大于0表示正在执行中
	concurrent bool                                              // 同一个任务是否允许并行执行
}

//通过任务去创建cron job(此处要区分运行文件(要指定路径)和一些指令)
func newJobFromTask(task *model.Task) (*Job, error) {
	job := newCommandJob(task)
	job.task = task
	job.concurrent = task.Concurrent == 1

	return job, nil
}

func newCommandJob(task *model.Task) *Job {
	job := &Job{
		id:   task.Id,
		name: task.Name,
	}

	//TaskType 0:Shell脚本, 1: 文件, 2:API
	job.runFunc = func(timeout time.Duration) (string, string, error, bool) {
		if task.TaskType == 1 {
			bufOut := new(bytes.Buffer)
			bufErr := new(bytes.Buffer)

			shellExt := model.LinuxShellExt
			if model.Common.SystemName == model.SystemWindows {
				shellExt = model.WindowsShellExt
			}
			runShell := fmt.Sprintf("%s/%s/%s.%s", model.WorkerRunDir, task.RunFilefolder, model.WorkerFileRunDir, shellExt)

			var cmd *exec.Cmd
			if model.Common.SystemName == model.SystemLinux {
				cmd = exec.Command("/bin/bash", "-c", runShell)
			} else {
				//windos下面不知怎么的，运行bat文件要绝对路径, 但bat文件中cd又可以写相对路径
				dir,_ := filepath.Abs(filepath.Dir(os.Args[0]))
				cmd = exec.Command("cmd.exe", "/c", fmt.Sprintf("%s/%s", dir,runShell))
			}

			cmd.Stdout = bufOut
			cmd.Stderr = bufErr
			cmd.Start()
			err, isTimeout := runCmdWithTimeOut(cmd, timeout)

			encoder := mahonia.NewDecoder("gbk")
			return encoder.ConvertString(bufOut.String()), encoder.ConvertString(bufErr.String()), err, isTimeout

		} else if task.TaskType == 2 {
			header := make(http.Header)
			if task.ApiHeader != "" && strings.TrimSpace(task.ApiHeader) != "" {
				headers := strings.Split(task.ApiHeader, "\n")
				for _,val := range headers {
					keyval := strings.Split(val, "=")
					if len(keyval) > 0 {
						v := strings.TrimSpace(keyval[0])
						v1 := strings.TrimSpace(keyval[1])
						if v != "" && v1 != "" {
							header.Set(v, v1)
						} else {
							continue
						}
					}
				}
			}
			responsestr := ""
			var err error
			var res *req.Resp

			req.SetTimeout(time.Second * time.Duration(task.TimeOut))
			if task.TaskApiMethod == "POST" {
				if task.ApiBody != "" {
					contenttype := header.Get("Content-Type")
					if contenttype == "" {
						header.Set("Content-Type","application/x-www-form-urlencoded")
					}
					if contenttype == "application/json" {
						res, err = req.Post(task.TaskApiUrl, header, req.BodyJSON(task.ApiBody))
					} else if contenttype == "application/xml" {
						res, err = req.Post(task.TaskApiUrl, header, req.BodyXML(task.ApiBody))
					} else  {
						res, err = req.Post(task.TaskApiUrl, header, task.ApiBody)
					}
				} else {
					res, err = req.Post(task.TaskApiUrl, header)
				}

			} else {
				res, err = req.Get(task.TaskApiUrl, header)
			}

			if err == nil {
				bodystr, _ := ioutil.ReadAll(res.Response().Body)
				defer res.Response().Body.Close()

				responsestr = string(bodystr)

				if res.Response().StatusCode != 200 {
					return responsestr, "", errors.New(fmt.Sprintf("返回的状态码为：%d", res.Response().StatusCode)), false
				}

				return responsestr, "", nil, false
			}
			return "", "", err, false

		} else if task.TaskType == 0 {
			bufOut := new(bytes.Buffer)
			bufErr := new(bytes.Buffer)

			var cmd *exec.Cmd
			if model.Common.SystemName == model.SystemLinux {
				cmd = exec.Command("/bin/bash", "-c", task.Command)
			} else {
				cmd = exec.Command("cmd.exe", "/c", task.Command)
			}

			cmd.Stdout = bufOut
			cmd.Stderr = bufErr
			cmd.Start()
			err, isTimeout := runCmdWithTimeOut(cmd, timeout)

			encoder := mahonia.NewDecoder("gbk")
			return encoder.ConvertString(bufOut.String()), encoder.ConvertString(bufErr.String()), err, isTimeout
		} else {
			return "", "", nil, false
		}
	}
	return job
}

func (this *Job) Run() {
	if !this.concurrent && this.status > 0 {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			//此处最好写日志
			fmt.Printf("Run wrong is : %s\n", err)
		}
	}()

	//此处是为了控制同时运行任务的个数
	if workPool != nil {
		workPool <- true
		defer func() {
			<-workPool
		}()
	}

	this.status++
	defer func() {
		this.status--
		if this.status < 0 {
			this.status = 0
		}
	}()

	t := time.Now()
	timeout := time.Duration(time.Hour * 24)
	if this.task.TimeOut > 0 {
		timeout = time.Second * time.Duration(this.task.TimeOut)
	}

	cmdOut, cmdErr, err, isTimeout := this.runFunc(timeout)
	ut := time.Now().Sub(t) / time.Millisecond

	//更新任务执行时间等
	if err := data.UpdateBackTask(t.Unix(),this.task.Id); err != nil {
		fmt.Printf("update task has error: %s", err.Error())
	}

	//写日志
	log := &model.TaskLog{
		TaskId 		: this.task.Id,
		Output 		: cmdOut,
		Error  		: cmdErr,
		ProcessTime : int(ut),
		CreateTime  : t.Unix(),
		Status      : model.TASK_SUCCESS,
	}
	if isTimeout {
		log.Status = model.TASK_TIMEOUT
		log.Error = fmt.Sprintf("任务执行超过 %d 秒\n---------------\n%s\n", int(timeout/time.Second), cmdErr)
	} else if err != nil {
		log.Status = model.TASK_ERROR
		log.Error = err.Error() + ":" + cmdErr
	}
	data.AddTaskLog(log)

	//发邮件

	/*
	if (this.task.Notify == 1 && err != nil) || this.task.Notify == 2 {

	}
	*/
}

func runCmdWithTimeOut(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		//超时，记日志等
		fmt.Printf("任务执行时间超过%d秒，进程将被强制杀掉: %d", int(timeout/time.Second), cmd.Process.Pid)
		go func() {
			<-done
		}()
		if err = cmd.Process.Kill(); err != nil {
			fmt.Printf("进程无法杀掉: %d, 错误信息: %s", cmd.Process.Pid, err)
		}
		return err, true

	case err = <-done:
		return err, false
	}
}
