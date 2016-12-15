package main

import (
	"fmt"
	"os"

	"github.com/mediocregopher/radix.v2/redis"
)

const RoomKey = "room"

func main() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer c.Close()

	c.Cmd("FLUSHALL")

	ids := []int{3, 2, 1}
	c.Cmd("LPUSH", RoomKey, ids)

	ids = []int{7, 8, 9}
	c.Cmd("RPUSH", RoomKey, ids)

	values, err := c.Cmd("LRANGE", RoomKey, 0, -1).Array()
	if err != nil {
		fmt.Println(err)
	}

	ids = []int{}
	for _, val := range values {
		id, _ := val.Int()
		ids = append(ids, id)
	}

	fmt.Println(ids)
	// Output:
	// [1 2 3 7 8 9]

	id, err := c.Cmd("RPOP", RoomKey).Int()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(id)
	// Output:
	// 9
}
