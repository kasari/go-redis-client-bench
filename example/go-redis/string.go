package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/redis.v5"
)

type User struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func main() {
	c := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	defer c.Close()

	c.FlushAll()

	kasari := User{1, "kasari"}
	b, _ := json.Marshal(kasari)

	c.Set(string(kasari.ID), b, 0)

	userJSON, err := c.Get(string(kasari.ID)).Bytes()
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

	c.Set(string(vongole.ID), b, 0)

	userJSONs, err := c.MGet(string(kasari.ID), string(vongole.ID)).Result()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var users []User
	for _, userJSON := range userJSONs {
		var u User
		json.Unmarshal([]byte(userJSON.(string)), &u)
		users = append(users, u)
	}

	fmt.Printf("%+v\n", users)
	// Output:
	// [{ID:1 Name:kasari} {ID:2 Name:vongole}]
}
