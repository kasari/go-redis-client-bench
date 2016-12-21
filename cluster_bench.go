package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	redigoCluster "github.com/chasex/redis-go-cluster"
	"github.com/kayac/parallel-benchmark/benchmark"
	radixCluster "github.com/mediocregopher/radix.v2/cluster"
	goredis "gopkg.in/redis.v5"
)

func newGoRedisCluster() *goredis.ClusterClient {
	ClusterNodes := []string{"127.0.0.1:7001", "127.0.0.1:7002"}
	cluster := goredis.NewClusterClient(
		&goredis.ClusterOptions{
			Addrs:       ClusterNodes,
			PoolSize:    10,
			IdleTimeout: 60 * time.Second,
		})
	return cluster
}

func newRedigoCluster() *redigoCluster.Cluster {
	ClusterNodes := []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"}
	cluster, err := redigoCluster.NewCluster(
		&redigoCluster.Options{
			StartNodes: ClusterNodes,
			KeepAlive:  10,
			AliveTime:  60 * time.Second,
		})
	if err != nil {
		log.Fatalf("redigo.New error: %s", err.Error())
		os.Exit(1)
	}

	return cluster
}

func newRadixCluster() *radixCluster.Cluster {
	cluster, err := radixCluster.NewWithOpts(
		radixCluster.Opts{
			Addr:         "127.0.0.1:7000",
			PoolSize:     10,
			PoolThrottle: 60 * time.Second,
		})
	if err != nil {
		log.Fatalf("radix.New error: %s", err.Error())
		os.Exit(1)
	}

	return cluster
}

func main() {
	var bytes, count, sec, parallel int
	var host, port, command string
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.StringVar(&port, "p", "7000", "port")
	flag.IntVar(&bytes, "b", 10, "byte")
	flag.IntVar(&count, "c", 5, "count")
	flag.IntVar(&parallel, "parallel", runtime.NumCPU(), "parallel")
	flag.IntVar(&sec, "s", 5, "sec")
	flag.StringVar(&command, "com", "SET", "bench command")
	flag.Parse()

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Println(addr)
	duration := time.Duration(sec) * time.Second
	fmt.Println(duration)

	switch command {
	case "SET":
		benchmarkSET(addr, bytes, count, parallel, duration)
		// case "GET":
		// 	benchmarkGET(addr, bytes, count, parallel, duration)
	}
}

func generateKey() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	runes := make([]rune, 4)
	for i := range runes {
		runes[i] = letters[rand.Intn(len(letters))]
	}

	return string(runes)
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
		fmt.Printf("%d\n", int(float64(score)/elapsed.Seconds()))
		// fmt.Printf("[%s] %d\n", clientName, int(float64(score)/elapsed.Seconds()))
	}
}

func benchmarkSET(addr string, bytes, count, parallel int, duration time.Duration) {
	resultMap := make(map[string][]*benchmark.Result)

	for ct := 1; ct <= count; ct++ {
		fmt.Printf("\n=================== SET (count: %d) ===================\n", ct)

		{ // GoRedis
			fmt.Println("[ GoRedis ]")

			c := newGoRedisCluster()

			result := benchmark.RunFunc(
				func() (subscore int) {
					if err := c.Set(generateKey(), strings.Repeat("a", bytes), 0).Err(); err != nil {
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

			c := newRedigoCluster()

			result := benchmark.RunFunc(
				func() (subscore int) {
					if _, err := c.Do("SET", generateKey(), strings.Repeat("a", bytes)); err != nil {
						return 0
					}
					return 1
				},
				duration,
				parallel,
			)
			resultMap["redigo"] = append(resultMap["redigo"], result)

			c.Close()
		}

		{ // Radix
			fmt.Println("[ Radix ]")

			c, err := radixCluster.New(addr)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			result := benchmark.RunFunc(
				func() (subscore int) {
					if val := c.Cmd("SET", generateKey(), strings.Repeat("a", bytes)); val.Err != nil {
						return 0
					}
					return 1
				},
				duration,
				parallel,
			)
			resultMap["radix"] = append(resultMap["radix"], result)

			c.Close()
		}
	}

	printResult(resultMap)
}

func benchmarkGET(addr string, bytes, count, parallel int, duration time.Duration) {
	resultMap := make(map[string][]*benchmark.Result)

	for ct := 1; ct <= count; ct++ {
		fmt.Printf("\n=================== SET (count: %d) ===================\n", ct)

		{ // GoRedis
			fmt.Println("[ GoRedis ]")

			c := newGoRedisCluster()

			result := benchmark.RunFunc(
				func() (subscore int) {
					if err := c.Set(generateKey(), strings.Repeat("a", bytes), 0).Err(); err != nil {
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

			c := newRedigoCluster()

			result := benchmark.RunFunc(
				func() (subscore int) {
					if _, err := c.Do("SET", generateKey(), strings.Repeat("a", bytes)); err != nil {
						return 0
					}
					return 1
				},
				duration,
				parallel,
			)
			resultMap["redigo"] = append(resultMap["redigo"], result)

			c.Close()
		}

		{ // Radix
			fmt.Println("[ Radix ]")

			c := newRadixCluster()

			result := benchmark.RunFunc(
				func() (subscore int) {
					if val := c.Cmd("SET", generateKey(), strings.Repeat("a", bytes)); val.Err != nil {
						return 0
					}
					return 1
				},
				duration,
				parallel,
			)
			resultMap["radix"] = append(resultMap["radix"], result)

			c.Close()
		}
	}

	printResult(resultMap)
}
