package main

import (
	//"fmt"
	"context"
	"time"
	"math/rand"
	//"strconv"
	//"ScheduleTask/test/rabbitmq"
	"fmt"
	"ScheduleTask/utils/system"
	"strings"
)

func main() {
	/*
		d := time.Now().Add(1 * time.Second)
		//50毫秒到了，触发如下代码
		ctx, cancel := context.WithDeadline(context.Background(), d)


		defer cancel()

		select {
		case <-time.After(3 * time.Second):
			fmt.Println("overslept")
		case <-ctx.Done():
			//50毫秒到了，执行该代码
			fmt.Println(ctx.Err())
		}
	*/

	/*
	 sliceT := []int{4,5,6}

	 slice1 := sliceT[0:1]
	 slice2 := sliceT[1:3]

	 fmt.Println("slice1:", slice1)
	 fmt.Println("slice2:", slice2[1])
	*/

	/*
		var msg chan string

		fmt.Println(msg)
		msg <- "dddd"
	*/

	//rabbitmq.Publish()
	//rabbitmq.SingleRecieve()
	//rabbitmq.RefleshConnection()
	//testslice()

	//testContext()

	key := system.CryptoSHA256("test")
	fmt.Println(key)

	var (
		err error
		token string
	)

	if token, err = system.Encrypt(`{"aaaa":"cccc","bb":456}`, "11"); err != nil {
		fmt.Println(err)
	} else {

	}

	fmt.Println(token)
	changeToken := strings.Replace(token, "9", "b", -1)
	fmt.Println(changeToken)

	value, flag := system.Decrypt(changeToken, "11")
	if flag {
		fmt.Println(value)
	} else {
		fmt.Print("dddd")
	}
}

func testContext()  {
	ctx, cancel := context.WithCancel(context.Background())
	eatNum := chiHanBao(ctx)
	for n := range eatNum {
		if n >= 10 {
			cancel()
			break
		}
	}

	fmt.Println("正在统计结果。。。")
	time.Sleep(2 * time.Second)
}

func chiHanBao(ctx context.Context) <-chan int {
	c := make(chan int)
	// 个数
	n := 0
	// 时间
	t := 0
	go func() {
		for {
			time.Sleep(time.Second)
			select {
			case <-ctx.Done():
				fmt.Printf("耗时 %d 秒，吃了 %d 个汉堡 \n", t, n)
				return
			case c <- n:
				incr := rand.Intn(5)
				n += incr
				if n >= 10 {
					n = 10
				}
				t++
				fmt.Printf("我吃了 %d 个汉堡\n", n)
			}
		}
	}()
	return c
}



func testslice() {
	var arr = [4]int{1, 2, 3, 4}

	slice1 := arr[1:3]
	fmt.Println(slice1)

	arr[1] = 6

	fmt.Println(slice1)

	slice1[0] = 8

	fmt.Println(arr)

	m := make(map[string]string, 10)

	m["aa"] = "cc"
	fmt.Println(m)
	fmt.Println(len(m))
}
