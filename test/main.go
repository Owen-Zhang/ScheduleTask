package main

import (
	//"fmt"
	//"context"
	//"time"
	//"strconv"
	//"ScheduleTask/test/rabbitmq"
	"fmt"
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

	err := fmt.Errorf("ddd%s", "ssss")
	fmt.Println(err)
	fmt.Println("dddd")

	testslice()
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
