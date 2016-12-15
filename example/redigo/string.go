package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
)

type User struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func main() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer c.Close()

	c.Do("FLUSHALL")

	kasari := User{1, "kasari"}
	b, _ := json.Marshal(kasari)

	c.Do("SET", kasari.ID, b)

	userJSON, err := redis.Bytes(c.Do("GET", kasari.ID))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var u User
	err = json.Unmarshal(userJSON, &u)

	fmt.Printf("%+v\n", u)
	// Output:
	// {ID:1 Name:kasari}

	// Multi Get
	vongole := User{2, "vongole"}
	b, _ = json.Marshal(vongole)

	c.Do("SET", vongole.ID, b)

	userJSONs, err := redis.ByteSlices(c.Do("MGET", kasari.ID, vongole.ID))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var users []User
	for _, userJSON := range userJSONs {
		var u User
		json.Unmarshal(userJSON, &u)
		users = append(users, u)
	}

	fmt.Printf("%+v\n", users)
	// Output:
	// [{ID:1 Name:kasari} {ID:2 Name:vongole}]
}
