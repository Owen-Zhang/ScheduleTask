package etc

import (
	"io/ioutil"
	"ScheduleTask/taskworker/api"
	"ScheduleTask/taskworker/jobs"
	workermodel "ScheduleTask/taskworker/model"
	"ScheduleTask/model"
	"ScheduleTask/storage"
	"gopkg.in/yaml.v2"
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
		Key  string `yaml:"key,omitempty"`
		Ip   string `yaml:"ip,omitempty"`
		Bind string `yaml:"bind,omitempty"`
		Note string `yaml:"note,omitempty"`
	} `yaml:"apiserver,omitempty"`

	CronService struct {
		PoolSize int32 `yaml:"poolsize,omitempty"`
	} `yaml:"cron,omitempty"`
	
	FileServer struct {
		Hosts   string `yaml:"hosts,omitempty"`
		Port 	int    `yaml:"port,omitempty"`
	} `yaml:"fileserver,omitempty"`

	CenterInfo struct {
		Hosts   string `yaml:"hosts,omitempty"`
		Port 	int    `yaml:"port,omitempty"`
	}

	WorkerInfo struct {
		Identification int `yaml:"identification,omitempty"`
	} `yaml:"workerinfo,omitempty"`
	
}

var configuration *Configuration

func newconfig(file string) error {
	data, err := ioutil.ReadFile(file)
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
			Key  string `yaml:"key,omitempty"`
			Ip   string `yaml:"ip,omitempty"`
			Bind string `yaml:"bind,omitempty"`
			Note string `yaml:"note,omitempty"`
		}{
			Key : "",
			Ip  : "192.168.0.103",
			Bind: ":8985",
			Note: "worker机器相关说明",
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

		CenterInfo: struct {
			Hosts  string `yaml:"hosts,omitempty"`
			Port   int    `yaml:"port,omitempty"`
		}{
			Hosts:  "127.0.0.1",
			Port: 8000,
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
			Key : configuration.ApiServer.Key,
			Ip  : configuration.ApiServer.Ip,
			Bind: configuration.ApiServer.Bind,
			Note: configuration.ApiServer.Note,
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

func GetCenterInfo() *workermodel.CenterInfo  {
	if configuration != nil {
		return &workermodel.CenterInfo{
			Hosts: configuration.CenterInfo.Hosts,
			Port : configuration.CenterInfo.Port,
		}
	}
	return nil
}
