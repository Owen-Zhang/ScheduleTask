package etc

import "ScheduleTask/storage"
import "gopkg.in/yaml.v2"

import (
	"io/ioutil"
	"ScheduleTask/taskworker/api"
	"os"
	"ScheduleTask/taskworker/jobs"
	"ScheduleTask/model"
)

type Configuration struct {
	Storage struct {
		Hosts  string `yaml:"hosts,omitempty"`
		DBName string `yaml:"dbname,omitempty"`
		User   string   `yaml:"user,omitempty"`
		Password string  `yaml:"password,omitempty"` 
		Port int  `yaml:"port,omitempty"`
	} `yaml:"storage,omitempty"`

	ApiServer struct {
		Bind string `yaml:"bind,omitempty"`
	} `yaml:"apiserver,omitempty"`

	CronService struct {
		PoolSize int32 `yaml:"poolsize,omitempty"`
	} `yaml:"cron,omitempty"`
	
	FileServer struct {
		Hosts   string `yaml:"hosts,omitempty"`
		Port 	int    `yaml:"port,omitempty"`
	} `yaml:"fileserver,omitempty"`
	
	WorkerInfo struct {
		Identification int `yaml:"identification,omitempty"`
	} `yaml:"workerinfo,omitempty"`
	
}

var configuration *Configuration

func New(file string) error {
	fp, err := os.OpenFile(file, os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}

	c := makeDefault()
	if err := yaml.Unmarshal([]byte(data), c); err != nil {
		return err
	}

	configuration = c
	return nil
}

func makeDefault() *Configuration {

	return &Configuration{
		Storage: struct {
			Hosts  string `yaml:"hosts,omitempty"`
			DBName string `yaml:"dbname,omitempty"`
			User   string   `yaml:"user,omitempty"`
			Password string  `yaml:"password,omitempty"` 
			Port int  `yaml:"port,omitempty"`
		}{
			Hosts:  "127.0.0.1",
			DBName: "jobschedule",
			User:   "guest",
			Password: "123456", 
			Port: 3306,
		},

		ApiServer: struct {
			Bind string `yaml:"bind,omitempty"`
		}{
			Bind: ":8985",
		},

		CronService: struct {
			PoolSize int32 `yaml:"poolsize,omitempty"`
		}{
			PoolSize: 10,
		},
		
		FileServer: struct {
			Hosts  string `yaml:"hosts,omitempty"`
			Port   int    `yaml:"port,omitempty"`
		}{
			Hosts:  "127.0.0.1",
			Port: 8988,
		},
		
		WorkerInfo: struct {
			Identification int `yaml:"identification,omitempty"`
		} {
			Identification: 1,
		},
	}
}

//数据访问的相关配制
func GetStorageArg() *storage.DataStorageArgs {

	if configuration != nil {
		return &storage.DataStorageArgs{
			Hosts:  configuration.Storage.Hosts,
			DBName: configuration.Storage.DBName,
			User:   configuration.Storage.User,
			Password: configuration.Storage.Password,
			Port:  configuration.Storage.Port,
		}
	}
	return nil
}

//对外api的相关配制
func GetApiServerArg() *api.ApiServerArg {
	if configuration != nil {
		return &api.ApiServerArg{
			Bind: configuration.ApiServer.Bind,
		}
	}
	return nil
}

//cron的相关配制
func GetCronArg() *jobs.CronArg {
	if configuration != nil {
		return &jobs.CronArg{
			PoolSize: configuration.CronService.PoolSize,
		}
	}
	return nil
}

//获取文件服务器相关的信息
func GetFileServerInfo() *model.FileServerInfo {
	if configuration != nil {
		return &model.FileServerInfo{
			Hosts: configuration.FileServer.Hosts,
			Port : configuration.FileServer.Port,
		}
	}
	return nil
}

func GetWorkerInfo() *model.WorkerInfo {
	if configuration != nil {
		return &model.WorkerInfo{
			Identification :configuration.WorkerInfo.Identification,
		}
	}
	return nil
}
