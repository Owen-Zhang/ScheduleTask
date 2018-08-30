package healthy

import (
	"fmt"
	"time"
	"sort"
	"strings"
	"net/http"
	"ScheduleTask/model"
	"github.com/imroc/req"
	"ScheduleTask/storage"
	"ScheduleTask/utils/system"
	"github.com/astaxie/beego/logs"
)

// 对实体进行排序
type byTime []*system.HealthInfo

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	if s[i] == nil {
		return false
	}

	if s[j] == nil {
		return true
	}

	return s[i].TimeOut.Before(s[j].TimeOut)
}

type healthy struct{
	 timeLocation  *time.Location
	 workerList    []*system.HealthInfo
	 add           chan *system.HealthInfo
}

var Health *healthy
var dataAccess *storage.DataStorage

func init()  {
	Health = &healthy{
		timeLocation:time.Now().Location(),
		workerList: make([]*system.HealthInfo, 0, 6),
		add: make(chan *system.HealthInfo),
	}
}

// 将报告数据加入到通道中
func Input(info *system.HealthInfo) {
	Health.add <- &system.HealthInfo{
		Status: 0,
		TimeOut:time.Now().Add(3 * time.Minute),//.In(Health.timeLocation),
		WorkerInfo:system.WorkerInfo{
			Name:info.WorkerInfo.Name,
			Ip: info.WorkerInfo.Ip,
			Port:info.WorkerInfo.Port,
			OsName:info.WorkerInfo.OsName,
			Note:info.WorkerInfo.Note,
		},
	}
}

// FindWorker 通过运行平台查找worker机器,返回ip, port信息
func FindWorker(system string) (string, string) {
	for _,val := range Health.workerList {
		if val.Status == -1 {
			continue;
		}
		if strings.EqualFold(val.WorkerInfo.OsName, system) {
			return val.WorkerInfo.Ip, val.WorkerInfo.Port
		}
	}
	return "",""
}

func add(info *system.HealthInfo)  {
	one := findOne(info.WorkerInfo.Ip, info.WorkerInfo.Port)
	if one == nil {
		Health.workerList = append(Health.workerList, info)
		worker := &system.WorkerInfo{
			Name:info.WorkerInfo.Name,
			Ip:info.WorkerInfo.Ip,
			Port:info.WorkerInfo.Port,
			OsName:info.WorkerInfo.OsName,
			Note: "worker正式向中心报告",
		}
		if err := dataAccess.AddWorkerlogs(worker, 1); err != nil {
			logs.Error("first AddWorkerlogs has wrong: %s", err.Error())
		}

	} else {
		one.Status = 0
		one.TimeOut = info.TimeOut
	}
}

func findOne(ip, port string) *system.HealthInfo {
	for index, val := range Health.workerList {
		if val.WorkerInfo.Ip == ip && val.WorkerInfo.Port == port && val.Status == 0 {
			return Health.workerList[index]
		}
	}
	return nil
}

// 检查客户端报告状态
func CheckWorkerStatus(access *storage.DataStorage) {
	dataAccess = access

	defer func() {
		if err := recover(); err != nil {
			logs.Error(fmt.Sprintf("心跳检查出错: %s", err))
		}
	}()

	now := time.Now()//.In(Health.timeLocation)
	for {
		sort.Sort(byTime(Health.workerList))

		var timer *time.Timer
		if len(Health.workerList) == 0 { //|| Health.workerList[0].TimeOut.Before(now) {
			//fmt.Println("0")
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			//fmt.Println("1")
			//fmt.Println(Health.workerList[0].TimeOut.Sub(now))
			timer = time.NewTimer(Health.workerList[0].TimeOut.Sub(now))
		}

		for {
			select {
			case now = <-timer.C:
				//now = now.In(Health.timeLocation)
				flagRemove := false
				for index, val := range Health.workerList {
					//fmt.Println("now: ", now.Format("2006-01-02 15:04:05"))
					//fmt.Println("timeout: ", val.TimeOut.Format("2006-01-02 15:04:05"))

					if val.TimeOut.After(now) || val.Status == -1 {
						break
					}
					flagRemove = true
					Health.workerList[index].Status = -1
				}

				if flagRemove {
					tempList := make([]*system.HealthInfo, 0, 6)
					for _,val := range Health.workerList {
						if val.Status == 0 {
							tempList = append(tempList, val)
						} else {
							go transferTask(val)
						}
					}
					Health.workerList = tempList
				}

			case newEntry := <- Health.add:
				now = time.Now()
				timer.Stop()
				add(newEntry)
			}
			break
		}
	}
}

// 转移此机器上的任务
func transferTask(info *system.HealthInfo) {
	oldIp, oldPort := info.WorkerInfo.Ip, info.WorkerInfo.Port
	workerInfo := fmt.Sprintf("%s_%s",oldIp, oldPort)

	worker := &system.WorkerInfo{
		Name:info.WorkerInfo.Name,
		Ip:info.WorkerInfo.Ip,
		Port:info.WorkerInfo.Port,
		OsName:info.WorkerInfo.OsName,
	}

	newIp, newPort := FindWorker(info.WorkerInfo.OsName);
	if newIp == "" || newPort == "" {
		if err := dataAccess.BatchUpdateTaskStatusByWorkerInfo(workerInfo, "", 0); err != nil {
			logs.Error("BatchUpdateTaskStatusByWorkerInfo has wrong: %s", err.Error())
		}
		worker.Note = fmt.Sprintf("此worker不能正常的向中心报告状态, 同时系统中又没有相同运行平台【%s】的woker, 请管理员处理", info.WorkerInfo.OsName)
		dataAccess.AddWorkerlogs(worker, 0)
		return
	}

	taskIds := dataAccess.GetTaskByWorkerInfo(workerInfo)
	if taskIds == nil {
		return
	}
	ids := strings.Join(taskIds, ",")
	if errupdate := dataAccess.BatchUpdateTaskStatusAndWorkerInfo(ids, 1, workerInfo); errupdate != nil {
		logs.Error("BatchUpdateTaskStatusAndWorkerInfo has wrong : %s", errupdate.Error())
	}

	worker.Note = fmt.Sprintf("此worker不能正常的向中心报告状态, 我们会将此worker上的任务转移到【%s】【%s】woker上", newIp, newPort)
	dataAccess.AddWorkerlogs(worker, 0)

	posturl := fmt.Sprintf(model.WorkerUrl, newIp, newPort, "batchstarttask")
	res, err := req.Post(posturl, req.Param{"ids" : ids})
	if err != nil || res.Response().StatusCode != http.StatusOK {
		if err == nil {
			logs.Error(
				"[转移worker]通知客戶端失敗, oldIp:%s, oldPort:%s; newIp:%s, newPort:%s;", oldIp,oldPort,newIp,newPort)
		} else {
			logs.Error(
				"[转移worker]通知客戶端失敗, oldIp:%s, oldPort:%s; newIp:%s, newPort:%s; error info:%s", oldIp,oldPort,newIp,newPort, err.Error())
		}
		if errupdate := dataAccess.BatchUpdateTaskStatusAndWorkerInfo(ids, 0, ""); errupdate != nil {
			logs.Error("BatchUpdateTaskStatusAndWorkerInfo has wrong : %s", errupdate.Error())
		}
	}
}
