package rabbitmq

import (
	"fmt"
	"strconv"
	"github.com/streadway/amqp"
)

// 发布信息
func Publish()  {
	var (
		err error
		con *amqp.Connection
		ch  *amqp.Channel
		q    amqp.Queue
	)
	con,err = amqp.Dial("amqp://admin:admin@192.168.0.109:5672/my_vhost")
	if err != nil {
		fmt.Println("connection: " + err.Error())
		return
	}
	defer con.Close()

	ch, err = con.Channel()
	if err != nil {
		fmt.Println("Channel: "+ err.Error())
		return
	}
	defer ch.Close()

	q, err = ch.QueueDeclare(
		"hello", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		fmt.Println("declare: "+ err.Error())
		return
	}


	for i := 0; i < 10; i++ {
		body := "third queue infomation: " + strconv.Itoa(i)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			fmt.Println("Publish: "+ err.Error())
		}
	}
	fmt.Println("send ok")
}