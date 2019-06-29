package rabbitmq

import (
	"fmt"
	"time"
	"github.com/streadway/amqp"
)

//同时订阅多个队列
func MultiReceive()  {
	var (
		err error
		con *amqp.Connection
		ch  *amqp.Channel
		q    amqp.Queue
		//msgs <- chan amqp.Delivery
		deliveryList []<- chan amqp.Delivery
		queues = []string{"hello","hello2","hello3"}
	)
	con,err = amqp.Dial("amqp://admin:admin@192.168.0.109:5672/my_vhost")
	if err != nil {
		fmt.Println("connection: " + err.Error())
		return
	}
	defer con.Close()


	for _,item := range queues {
		ch, err = con.Channel()
		if err != nil {
			fmt.Println("Channel: "+ err.Error())
			continue
		}
		defer ch.Close()

		q, err = ch.QueueDeclare(
			item, // name
			true,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)

		if err != nil {
			fmt.Println("declare: "+ err.Error())
			continue
		}
		var msgs <- chan amqp.Delivery

		msgs, err = ch.Consume(
			q.Name, // queue
			"",     // consumer
			false,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			fmt.Println("Consume: " + err.Error())
			continue
		}
		deliveryList = append(deliveryList, msgs)
	}


	forever := make(chan bool)

	go func(){
		for _, result := range deliveryList {
			fmt.Println(len(result))
			for msg := range result {
				go func(){
					time.Sleep(2000)
					fmt.Println(string(msg.Body))
					msg.Ack(true)
				}()
			}
		}
	}()

	<-forever
}