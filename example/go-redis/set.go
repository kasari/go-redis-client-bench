package main

import (
	"fmt"
	"os"
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

	ids := []int{1, 3, 5, 7, 9}
	var args []interface{}
	for _, id := range ids {
		args = append(args, id)
	}
	c.SAdd(RoomKey, args...)

	ids = []int{}
	values, err := c.SRandMemberN(RoomKey, 3).Result()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, val := range values {
		id, _ := strconv.Atoi(val)
		ids = append(ids, id)
	}

	fmt.Println(ids)
}
