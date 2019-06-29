package rabbitmq

import (
	"github.com/streadway/amqp"
	"errors"
	"fmt"
	"time"
)

var Connection *amqp.Connection

//重新连接
func RefleshConnection()  {
	var (
		err error
		//ch  *amqp.Channel
	)
	if Connection != nil {
		_,err = Connection.Channel()
	} else {
		err = errors.New("connection is nil")
	}

	if err != nil {
		for {
			Connection, err = amqp.Dial("amqp://admin:admin@192.168.0.109:5672/my_vhost")
			if err != nil {
				fmt.Println("after 10s connect")
				time.Sleep(10 * time.Second)
			} else {
				fmt.Println("connected !!!!!!")
				break;
			}
		}
	} else {
		fmt.Println("connection is alive")
	}
}