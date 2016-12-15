package main

import (
	"fmt"
	"strconv"

	"gopkg.in/redis.v5"
)

const RoomKey = "room"

func main() {
	c := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	defer c.Close()

	c.FlushAll()

	ids := []int{3, 2, 1}
	var args []interface{}
	for _, id := range ids {
		args = append(args, id)
	}
	c.LPush(RoomKey, args...)

	ids = []int{7, 8, 9}
	args = []interface{}{}
	for _, id := range ids {
		args = append(args, id)
	}
	c.RPush(RoomKey, args...)

	values, err := c.LRange(RoomKey, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}

	ids = []int{}
	for _, val := range values {
		id, _ := strconv.Atoi(val)
		ids = append(ids, id)
	}

	fmt.Println(ids)
	// Output:
	// [1 2 3 7 8 9]

	id, err := c.RPop(RoomKey).Int64()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(id)
	// Output:
	// 9
}
