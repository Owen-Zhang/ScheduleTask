package model

//worker相关信息
type Worker struct {
	Id     int

	/*名称*/
	Name   string

	/*worker机器必须要有此，不然不会加入可使用worker队列*/
	Key    string

	/*状态，服务端控制是否使用此机器, 0未加入队列, 1 加入队列中*/
	Status int8

	/*备注*/
	Note   string
}
