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

	ids := []int{1, 3, 5, 7, 9}
	c.Cmd("SADD", RoomKey, ids)

	values, err := c.Cmd("SRANDMEMBER", RoomKey, 3).Array()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ids = []int{}
	for _, val := range values {
		id, _ := val.Int()
		ids = append(ids, id)
	}

	fmt.Println(ids)
}
