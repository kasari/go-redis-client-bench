package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/kayac/parallel-benchmark/benchmark"
	radix "github.com/mediocregopher/radix.v2/redis"
	goredis "gopkg.in/redis.v5"
)

func newGoRedisClient(addr string) *goredis.Client {
	c := goredis.NewClient(&goredis.Options{
		Addr: addr,
	})

	return c
}

func newRedigoClient(addr string) redigo.Conn {
	c, err := redigo.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return c
}

func newRadixClient(addr string) *radix.Client {
	c, err := radix.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return c
}

func main() {
	var bytes, count, sec int
	var host, port, command string
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.StringVar(&port, "p", "6379", "port")
	flag.IntVar(&bytes, "b", 10, "byte")
	flag.IntVar(&count, "c", 1, "count")
	flag.IntVar(&sec, "s", 5, "sec")
	flag.StringVar(&command, "com", "SET", "bench command")
	flag.Parse()

	addr := fmt.Sprintf("%s:%s", host, port)
	duration := time.Duration(sec) * time.Second

	switch command {
	case "SET":
		benchmarkSET(addr, bytes, count, duration)
	case "GET":
		benchmarkGET(addr, bytes, count, duration)
	}
}

func printResult(resultMap map[string][]*benchmark.Result) {
	fmt.Println("[ Result ]")
	for _, clientName := range []string{"goredis", "redigo", "radix"} {
		var score int
		var elapsed time.Duration
		for _, r := range resultMap[clientName] {
			score += r.Score
			elapsed += r.Elapsed
		}
		fmt.Printf("[%s] %d\n", clientName, int(float64(score)/elapsed.Seconds()))
	}
}

func benchmarkSET(addr string, bytes, count int, duration time.Duration) {
	resultMap := make(map[string][]*benchmark.Result)

	for ct := 1; ct <= count; ct++ {
		fmt.Printf("\n=================== SET (count: %d) ===================\n", ct)

		{ // GoRedis
			fmt.Println("[ GoRedis ]")

			c := newGoRedisClient(addr)
			c.FlushAll()

			result := benchmark.RunFunc(
				func() (subscore int) {
					c.Set("hoge", strings.Repeat("a", bytes), 0)
					return 1
				},
				duration,
				1,
			)
			resultMap["goredis"] = append(resultMap["goredis"], result)
		}

		{ // Redigo
			fmt.Println("[ Redigo ]")

			c := newRedigoClient(addr)
			c.Do("flushall")

			result := benchmark.RunFunc(
				func() (subscore int) {
					c.Do("SET", "hoge", strings.Repeat("a", bytes))
					return 1
				},
				duration,
				1,
			)
			resultMap["redigo"] = append(resultMap["redigo"], result)
		}

		{ // Radix
			fmt.Println("[ Radix ]")

			c := newRadixClient(addr)
			c.Cmd("flushall")

			result := benchmark.RunFunc(
				func() (subscore int) {
					c.Cmd("SET", "hoge", strings.Repeat("a", bytes))
					return 1
				},
				duration,
				1,
			)
			resultMap["radix"] = append(resultMap["radix"], result)
		}
	}

	printResult(resultMap)
}

func benchmarkGET(addr string, bytes, count int, duration time.Duration) {
	resultMap := make(map[string][]*benchmark.Result)

	c := newGoRedisClient(addr)
	c.Set("hoge", strings.Repeat("a", bytes), 0)
	c.Close()

	for ct := 1; ct <= count; ct++ {

		fmt.Printf("\n=================== GET (count: %d) ===================\n", ct)

		{ // GoRedis
			fmt.Println("[ GoRedis ]")

			c := newGoRedisClient(addr)

			result := benchmark.RunFunc(
				func() (subscore int) {
					c.Get("hoge")
					return 1
				},
				duration,
				1,
			)
			resultMap["goredis"] = append(resultMap["goredis"], result)
		}

		{ // Redigo
			fmt.Println("[ Redigo ]")

			c := newRedigoClient(addr)

			result := benchmark.RunFunc(
				func() (subscore int) {
					c.Do("GET", "hoge")
					return 1
				},
				duration,
				1,
			)
			resultMap["redigo"] = append(resultMap["redigo"], result)
		}

		{ // Radix
			fmt.Println("[ Radix ]")

			c := newRadixClient(addr)

			result := benchmark.RunFunc(
				func() (subscore int) {
					c.Cmd("GET", "hoge")
					return 1
				},
				duration,
				1,
			)
			resultMap["radix"] = append(resultMap["radix"], result)
		}
	}

	printResult(resultMap)
}
