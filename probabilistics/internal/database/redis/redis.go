package myredis

import (
	"context"
	"fmt"
	"gopds/probabilistics/internal/config"
	"sync"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Ctx = context.Background()

var Initialize = sync.OnceFunc(func() {
	var err error
	config.LoadRedisConfig()
	opts := &redis.Options{
		Addr: config.RedisCfg().Addr,
	}
	Client = redis.NewClient(opts)

	err = Client.Ping(Ctx).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to Redis!")

})
