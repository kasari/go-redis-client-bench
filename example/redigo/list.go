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

	ids := []int{3, 2, 1}
	c.Do("LPUSH", redis.Args{}.Add(RoomKey).AddFlat(ids)...)

	ids = []int{7, 8, 9}
	c.Do("RPUSH", redis.Args{}.Add(RoomKey).AddFlat(ids)...)

	ids, err = redis.Ints(c.Do("LRANGE", RoomKey, 0, -1))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ids)
	// Output:
	// [1 2 3 7 8 9]

	id, err := redis.Int(c.Do("RPOP", RoomKey))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(id)
	// Output:
	// 9
}
