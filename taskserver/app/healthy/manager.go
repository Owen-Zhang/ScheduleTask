package healthy

import (
	"time"
	"fmt"
	"ScheduleTask/model"
	"net/http"
	"github.com/astaxie/beego"
	"github.com/imroc/req"
)

type Health struct {
	id 		   int                              // 任务ID
	name       string                           // 任务名称
	health     *model.HealthInfo                // 心跳对象信息
	runFunc    func(time.Duration) bool			// 执行函数
	count      int                              // check faild次数(等于三次就表示worker不能用了，报告状态，启用其它的机器)
	faild      bool								// 标识某个worker不可用
}

//新增心跳任务
func newHealth(heart *model.HealthInfo) (*Health, error) {
	hea := &Health{
		id:     heart.Id,
		name:   heart.Name,
		count:  0,
		faild:  false,
		health: heart,
	}
	hea.runFunc = func(timeout time.Duration) bool {
		return runWithTimeOut(heart.Url, heart.Port, timeout)
	}
	return hea, nil
}

//客户端要向服务器报告机子已好, 客户端和服务器双向交互，将faild状态改为true(现在要手动去服务端改状态)
func (this *Health) Run() {
	if this.faild {
		return
	}

	timeout := time.Duration(time.Second * 5)
	flag := this.runFunc(timeout)

	if !flag {
		this.count++
		fmt.Printf("check [%s:%d] faild", this.health.Url, this.health.Port)
	} else {
		this.count--
	}

	//此worker不可用, 提示相关信息，并将任务分配给其它的worker
	pingcount, err := beego.AppConfig.Int("worker.pingcount")
	if err != nil {
		fmt.Printf("[获取worker.pingcount配制：%s]", err)
	}

	if this.count >= pingcount {
		fmt.Printf("[%s:%d] 任务将要转移到其它worker", this.health.Url, this.health.Port)
		
		//看看当前出问题的机子有没有挂着的任务
		workerList, count := dataaccess.TaskGetList(1, 500, -1, 0, this.health.Id)
		if workerList == nil || count == 0 {
			return			
		}
		if workerList != nil && len(workerList) > 0 {
			//1: 处理转移任务
			list, errworker := dataaccess.GetWorkerList(1, this.health.SystemInfo)
			if errworker != nil {
				fmt.Println(errworker)
				return
			}
			if list != nil && len(list) > 0 {
				newworkerid := 0
				temphealth := &model.HealthInfo{}
				
				//取第一个作为转入的机器(没有考虑机器的性能，处理能力和业务数据等)
				for	_,val := range list {
					if val.Id != this.health.Id {
						newworkerid = val.Id
						temphealth = val
						break
					}
				}
				if newworkerid != 0 {
					err := dataaccess.UpdateTaskWorker(this.health.Id, newworkerid)
					if err != nil {
						//调用api，提交任务给另外的机子 //这里最好批量去处理
						posturl := fmt.Sprintf(model.WorkerUrl, temphealth.Url, temphealth.Port, "starttask")
						
						for _, val := range workerList {					
							res, err :=
							req.Post(posturl, req.Param{"id" : val.Id})
							
							if err != nil || res.Response().StatusCode != http.StatusOK {
								if err == nil {
									fmt.Println("[任务转移]通知客戶端失敗")
									//此处最好邮件什么的(或者记录特殊日志)
								}
							}
						}
					}
				}
			}
		}
		//2: 标识此worker不可用
		this.faild = true
		dataaccess.DeleteWorker(this.health.Id)
	}
}

func runWithTimeOut(url string, port int, timeout time.Duration) bool {
	client := http.Client{}
	client.Timeout = timeout
	respose, err := client.Post(fmt.Sprintf("http://%s:%d/worker/ping", url, port), "", nil)
	if err != nil {
		return false
	}

	return respose.StatusCode == http.StatusOK
}
