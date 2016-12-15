package main

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
)

const RoomKey = "room"

func main() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer c.Close()

	c.Do("FLUSHALL")

	ids := []int{1, 3, 5, 7, 9}
	c.Do("SADD", redis.Args{}.Add(RoomKey).AddFlat(ids)...)

	// memo Int64sは存在しない
	ids, err = redis.Ints(c.Do("SRANDMEMBER", RoomKey, 3))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(ids)
}
