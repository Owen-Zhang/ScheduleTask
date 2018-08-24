package healthy

import (
	"time"
	"sort"
	"fmt"
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

	return s[i].ReceiveTime.Before(s[j].ReceiveTime)
}

type healthy struct{
	 timeLocation  *time.Location
	 workerList    []*system.HealthInfo
	 add           chan *system.HealthInfo
}

var Health *healthy

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
		ReceiveTime:time.Now().In(Health.timeLocation),
		WorkerInfo:system.WorkerInfo{
			Name:info.WorkerInfo.Name,
			Ip: info.WorkerInfo.Ip,
			Port:info.WorkerInfo.Port,
			OsName:info.WorkerInfo.OsName,
			Note:info.WorkerInfo.Note,
		},
	}
}

func add(info *system.HealthInfo)  {
	one := findOne(info.WorkerInfo.Ip, info.WorkerInfo.Port)
	if one == nil {
		Health.workerList = append(Health.workerList, info)
	} else {
		one.Status = 0
		one.ReceiveTime = time.Now().In(Health.timeLocation)
	}
}

func findOne(ip, port string) *system.HealthInfo {
	for index, val := range Health.workerList {
		if val.WorkerInfo.Ip == ip && val.WorkerInfo.Port == port {
			return Health.workerList[index]
		}
	}
	return nil
}

// 检查客户端报告状态
func CheckWorkerStatus() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(fmt.Sprintf("心跳检查出错: %s", err))
		}
	}()

	now := time.Now().In(Health.timeLocation)
	for {
		sort.Sort(byTime(Health.workerList))

		var timer *time.Timer
		if len(Health.workerList) == 0{
			fmt.Println("0")
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			fmt.Println("1")
			timer = time.NewTimer(now.Add(3 * time.Minute).Sub(Health.workerList[0].ReceiveTime))
		}

		fmt.Println("ggggg")

		for {
			select {
			case now = <-timer.C:
				now = now.In(Health.timeLocation)
				for index, val := range Health.workerList {
					//fmt.Printf("receive + 3 minutes: %s", val.ReceiveTime.Add(3 * time.Minute).In(Health.timeLocation).Format("2006-01-02 15:04:05"))
					//fmt.Printf("now: %s", now.Format("2006-01-02 15:04:05"))
					fmt.Println(val.ReceiveTime.Add(3 * time.Minute).In(Health.timeLocation).Sub(now))

					if val.ReceiveTime.Add(3 * time.Minute).In(Health.timeLocation).Before(now) && val.Status != -1 {
						Health.workerList[index].Status = -1
						fmt.Println(Health.workerList[index].Status)
					}
				}
			case newEntry := <- Health.add:
				timer.Stop()
				add(newEntry)
			}
			break
		}
	}
}
