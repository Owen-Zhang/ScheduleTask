package healthy

import (
	"time"
	"ScheduleTask/utils/system"
	"sort"
)

// 对实体进行排序
type byTime []*system.HealthInfo

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	return s[i].ReceiveTime.Before(s[j].ReceiveTime)
}

type healthy struct{
	 timeLocation  *time.Location
	 workerList    []*system.HealthInfo
	 add           chan *system.HealthInfo
}

var Health healthy

func init()  {
	Health.timeLocation = time.Now().Location()
	Health.workerList = make([]*system.HealthInfo, 6)
	Health.add = make(chan *system.HealthInfo)
}

// 将报告数据加入到通道中
func Input(info *system.HealthInfo) {
	Health.add <- &system.HealthInfo{
		ReceiveTime:time.Now(),
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
		one.ReceiveTime = time.Now()
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
	//获取本地当前时间
	now := time.Now().In(Health.timeLocation)
	for {
		sort.Sort(byTime(Health.workerList))

		var timer *time.Timer

		if len(Health.workerList) == 0{
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			timer = time.NewTimer(now.Add(3 * time.Minute).Sub(Health.workerList[0].ReceiveTime))
		}

		for {
			select {
			case now = <-timer.C:
				now = now.In(Health.timeLocation)
				for index, val := range Health.workerList {
					if val.ReceiveTime.Add(3 * time.Minute).Before(now) && val.Status != -1 {
						Health.workerList[index].Status = -1
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
