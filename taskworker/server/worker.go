package server

import (
	"flag"
	"ScheduleTask/taskworker/api"
	"ScheduleTask/taskworker/ctrl"
	"ScheduleTask/taskworker/etc"
	"ScheduleTask/taskworker/jobs"
	"ScheduleTask/storage"
	"errors"
	"ScheduleTask/taskworker/global"
)

type JobWork struct {
	Controller *ctrl.Controller
	Storage    *storage.DataStorage
	Api        *api.ApiServer
}

func NewWorker() (*JobWork, error) {

	var etcfile string
	flag.StringVar(&etcfile, "f", "etc/worker.yml", "worker etc file.")
	flag.Parse()
	if err := etc.New(etcfile); err != nil {
		return nil, err
	}

	storagearg := etc.GetStorageArg()
	dataaccess, err := storage.NewDataStorage(storagearg)
	if err != nil {
		return nil, err
	}

	fileserverinfo := etc.GetFileServerInfo()
	if fileserverinfo == nil {
		return nil, errors.New("get FileServerInfo is wrong")
	}
	
	workerinfo := etc.GetWorkerInfo()
	if workerinfo == nil {
		return nil, errors.New("GetWorkerInfo is wrong")
	}
	
	controller := ctrl.NewController(dataaccess, fileserverinfo)
	apiserver := api.NewAPiServer(etc.GetApiServerArg(), controller)
	jobs.NewCron(etc.GetCronArg(), dataaccess)

	job := &JobWork{
		Controller: controller,
		Storage:    dataaccess,
		Api:        apiserver,
	}

	//加载本地的任务到任务队列中
	job.Controller.AddAutoRunSelfTask(workerinfo.Identification)

	/*worker的相关信息*/
	global.WorkerInformation = workerinfo;
	
	return job, nil
}

func (s *JobWork) Start() error {
	go s.Controller.ListenTask()
	s.Api.StartUp()

	return nil
}

func (s *JobWork) Stop() error {
	defer func() {
		s.Storage.Close()
		s.Controller.Close()
	}()
	return nil
}
