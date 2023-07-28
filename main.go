package main

import (
	"ReviveProject/src"
	"flag"
	"fmt"
	"istio.io/pkg/cache"
	"os"
	"strconv"
	"time"
)

func main() {
	port := flag.Int("port", 8000, "Port to listen to")
	key := flag.String("key", "", "Hypixel API-Key")
	cacheTime := flag.Int("cache", 60, "Time in seconds for global memory cache to speedup repeated requests")
	flag.Parse()

	if *port <= 0 || *port > 65536 {
		fmt.Println("port must be 0-65536. Given " + strconv.Itoa(*port))
		flag.PrintDefaults()
		return
	}
	if !(len(*key) >= 32 && len(*key) <= 36) {
		fmt.Println("key must be a valid UUID4. Given " + *key)
		flag.PrintDefaults()
		return
	}
	if *cacheTime < 0 {
		fmt.Println("cache must be equal or bigger than 0. Given " + strconv.Itoa(*cacheTime))
		flag.PrintDefaults()
		return
	}

	start(*port, *key, *cacheTime)
}

func start(port int, key string, cacheTime int) {
	expiringCache := cache.NewLRU(time.Duration(cacheTime)*time.Second,
		time.Duration(cacheTime/2)*time.Second,
		500)

	hypixel := src.NewHypixelApi(key)
	manager, err := src.NewEndPointManager("endpoints.json")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	src.Listen(*manager, hypixel, expiringCache, port)
}
