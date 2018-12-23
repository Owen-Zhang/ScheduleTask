package api

import "ScheduleTask/taskworker/ctrl"
import "github.com/gin-gonic/gin"

type ApiServerArg struct {
	Key  string
	Ip   string
	Bind string
	Note string
}

type ApiServer struct {
	bind       string
	controller *ctrl.Controller
	s          *gin.Engine
}

func NewAPiServer(arg *ApiServerArg, contr *ctrl.Controller) *ApiServer {
	server := gin.Default()
	apiserver := &ApiServer{
		bind:       arg.Bind,
		controller: contr,
		s:          server,
	}
	apiserver.InitRoute()

	return apiserver
}

func (this *ApiServer) StartUp() {
	this.s.Run(this.bind)
}
