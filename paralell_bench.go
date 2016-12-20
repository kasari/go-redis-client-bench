package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/kayac/parallel-benchmark/benchmark"
	"github.com/mediocregopher/radix.v2/pool"
	goredis "gopkg.in/redis.v5"
)

func newGoRedisClient(addr string) *goredis.Client {
	c := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		PoolSize: 10,
	})

	return c
}

func newRedigoPool(addr string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redigo.Conn, error) { return redigo.Dial("tcp", addr) },
	}
}

func main() {
	var bytes, count, sec, parallel int
	var host, port, command string
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.StringVar(&port, "p", "6379", "port")
	flag.IntVar(&bytes, "b", 10, "byte")
	flag.IntVar(&count, "c", 5, "count")
	flag.IntVar(&parallel, "parallel", runtime.NumCPU(), "parallel")
	flag.IntVar(&sec, "s", 5, "sec")
	flag.StringVar(&command, "com", "SET", "bench command")
	flag.Parse()

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Println(addr)
	duration := time.Duration(sec) * time.Second

	switch command {
	case "SET":
		benchmarkSET(addr, bytes, count, parallel, duration)
	case "GET":
		benchmarkGET(addr, bytes, count, parallel, duration)
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

func benchmarkSET(addr string, bytes, count, parallel int, duration time.Duration) {
	resultMap := make(map[string][]*benchmark.Result)

	for ct := 1; ct <= count; ct++ {
		fmt.Printf("\n=================== SET (count: %d) ===================\n", ct)

		{ // GoRedis
			fmt.Println("[ GoRedis ]")

			c := newGoRedisClient(addr)
			c.FlushAll()

			result := benchmark.RunFunc(
				func() (subscore int) {
					if err := c.Set("hoge", strings.Repeat("a", bytes), 0).Err(); err != nil {
						return 0
					}
					return 1
				},
				duration,
				parallel,
			)
			resultMap["goredis"] = append(resultMap["goredis"], result)

			c.Close()
		}

		{ // Redigo
			fmt.Println("[ Redigo ]")

			pool := newRedigoPool(addr)
			c := pool.Get()
			c.Do("flushall")
			c.Close()

			result := benchmark.RunFunc(
				func() (subscore int) {
					c := pool.Get()
					defer c.Close()
					if _, err := c.Do("SET", "hoge", strings.Repeat("a", bytes)); err != nil {
						return 0
					}
					return 1
				},
				duration,
				parallel,
			)
			resultMap["redigo"] = append(resultMap["redigo"], result)

			pool.Close()
		}

		{ // Radix
			fmt.Println("[ Radix ]")

			pool, err := pool.New("tcp", addr, 4)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			c, _ := pool.Get()
			c.Cmd("flushall")
			pool.Put(c)

			result := benchmark.RunFunc(
				func() (subscore int) {
					if val := pool.Cmd("SET", "hoge", strings.Repeat("a", bytes)); val.Err != nil {
						return 0
					}
					return 1
				},
				duration,
				parallel,
			)
			resultMap["radix"] = append(resultMap["radix"], result)

			pool.Empty()
		}
	}

	printResult(resultMap)
}

func benchmarkGET(addr string, bytes, count, parallel int, duration time.Duration) {
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
					if err := c.Get("hoge").Err(); err != nil {
						return 0
					}
					return 1
				},
				duration,
				parallel,
			)
			resultMap["goredis"] = append(resultMap["goredis"], result)

			c.Close()
		}

		{ // Redigo
			fmt.Println("[ Redigo ]")

			pool := newRedigoPool(addr)

			result := benchmark.RunFunc(
				func() (subscore int) {
					c := pool.Get()
					defer c.Close()

					if _, err := c.Do("GET", "hoge"); err != nil {
						return 0
					}
					return 1
				},
				duration,
				parallel,
			)
			resultMap["redigo"] = append(resultMap["redigo"], result)

			pool.Close()
		}

		{ // Radix
			fmt.Println("[ Radix ]")

			pool, err := pool.New("tcp", addr, 10)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			result := benchmark.RunFunc(
				func() (subscore int) {
					if val := pool.Cmd("GET", "hoge"); val.Err != nil {
						return 0
					}
					return 1
				},
				duration,
				parallel,
			)
			resultMap["radix"] = append(resultMap["radix"], result)

			pool.Empty()
		}
	}

	printResult(resultMap)
}
