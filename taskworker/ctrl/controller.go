package ctrl

import (
	"ScheduleTask/storage"
	"time"
	"ScheduleTask/model"
)

type Controller struct {
	Ticker     *time.Ticker
	Actionlist chan Action
	Storage    *storage.DataStorage
	FileServer *model.FileServerInfo
}

//外部接口传入的任务实体(chan实体)
type Action struct {
	ActionType int    	//操作类型
	Id         int 		//任务的主键
}


func NewController(storage *storage.DataStorage, fileserver *model.FileServerInfo) *Controller {
	list := make(chan Action, 10)
	return &Controller{
		Storage:    storage,
		Actionlist: list,
		FileServer:  fileserver,
	}
}

func (this *Controller) ListenTask() {
NEW_TICK_DURATION:
	this.Ticker = time.NewTicker(time.Second * 1)
	for {
		select {
		case newtask := <-this.Actionlist:
			this.Ticker.Stop()

			actiontype := newtask.ActionType
			switch actiontype {
			case 1,2:
				this.start(&newtask)
			case 3:
				this.stop(newtask.Id)
			case 4:
				this.delete(newtask.Id)
			case 5:
				this.run(newtask.Id)
			}

			goto NEW_TICK_DURATION
		}
	}
}

func (this *Controller) Close() {
	close(this.Actionlist)
}
