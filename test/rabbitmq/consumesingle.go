package rabbitmq

import (
	"fmt"
	"time"
	"github.com/streadway/amqp"
	//"reflect"
	"reflect"
)

//订阅单个队列
func SingleRecieve()  {
	var (
		err error
		con *amqp.Connection
		ch  *amqp.Channel
		q    amqp.Queue
		msgs <- chan amqp.Delivery
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

	err = ch.Qos(1, 0,false)
	if err != nil {
		fmt.Println("Qos: "+ err.Error())
		return
	}

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
		return
	}

	forever := make(chan bool)

	go func(){
		for msg := range msgs {
			go func(){
				time.Sleep(10000)
				fmt.Println(string(msg.Body))
				//fmt.Println(reflect.TypeOf(map[string]interface{}(msg.Headers)))
				fmt.Println(map[string]interface{}(msg.Headers)["test"])

				temp := map[string]interface{}(msg.Headers)["test"];

				fmt.Println(reflect.TypeOf(temp))
				value,flag := temp.(int64)
				fmt.Println(flag)
				fmt.Println(value)

				if flag && value != 10 {
					msg.Ack(false)
				} else {
					fmt.Println("header not contain")
				}
			}()
		}
	}()

	<-forever
}