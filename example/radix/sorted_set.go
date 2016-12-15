package main

import (
	"fmt"
	"os"

	"github.com/mediocregopher/radix.v2/redis"
)

type Player struct {
	ID    int
	Score int
}

const RankingKey = "ranking"

func main() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer c.Close()

	c.Cmd("FLUSHALL")

	players := []Player{
		Player{1, 100000},
		Player{2, 1000000},
		Player{3, 10},
		Player{4, 10000},
		Player{5, 1000},
	}
	var args []interface{}
	for _, p := range players {
		args = append(args, p.Score, p.ID)
	}

	c.Cmd("ZADD", RankingKey, args)

	// memo Int64sは存在しない
	values, err := c.Cmd("ZREVRANGE", RankingKey, 0, 2).Array()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var playerIDs []int
	for _, val := range values {
		id, _ := val.Int()
		playerIDs = append(playerIDs, id)
	}

	fmt.Println(playerIDs)
	// Output:
	// [2 1 4]
}
