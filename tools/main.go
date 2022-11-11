package main

import (
	"bytes"
	"github.com/bxcodec/httpcache"
	"github.com/bxcodec/httpcache/cache/redis"
	"net/http"
	"time"
)

func main() {

	// Run redis before running this:
	// docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest

	client := http.Client{}
	httpcache.NewWithRedisCache(&client, true, &redis.CacheOptions{
		Addr: "localhost:6379",
	}, time.Second*60)

	var b []byte
	request, err := http.NewRequest("GET", "https://example.com", bytes.NewReader(b))

	if err != nil {
		panic(err)
	}

	_, err = client.Do(request)

	if err != nil {
		panic(err)
	}

}
