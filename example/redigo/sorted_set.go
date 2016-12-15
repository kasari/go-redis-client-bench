package main

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
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

	c.Do("FLUSHALL")

	players := []Player{
		Player{1, 100000},
		Player{2, 1000000},
		Player{3, 10},
		Player{4, 10000},
		Player{5, 1000},
	}
	for _, player := range players {
		c.Do("ZADD", RankingKey, player.Score, player.ID)
	}

	// memo Int64sは存在しない
	playerIDs, err := redis.Ints(c.Do("ZREVRANGE", RankingKey, 0, 2))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(playerIDs)
	// Output:
	// [2 1 4]
}
