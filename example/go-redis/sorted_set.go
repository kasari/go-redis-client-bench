package main

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/redis.v5"
)

type Player struct {
	ID    int
	Score int
}

const RankingKey = "ranking"

func main() {
	c := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	defer c.Close()

	c.FlushAll()

	players := []Player{
		Player{1, 100000},
		Player{2, 1000000},
		Player{3, 10},
		Player{4, 10000},
		Player{5, 1000},
	}

	var zs []redis.Z
	for _, player := range players {
		zs = append(zs, redis.Z{
			Score:  float64(player.Score),
			Member: player.ID,
		})
	}

	c.ZAdd(RankingKey, zs...)
	values, err := c.ZRevRange(RankingKey, 0, 2).Result()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var playerIDs []int
	for _, val := range values {
		id, _ := strconv.Atoi(val)
		playerIDs = append(playerIDs, id)
	}

	fmt.Println(playerIDs)
	// Output:
	// [2 1 4]
}
