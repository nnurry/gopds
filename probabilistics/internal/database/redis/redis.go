package myredis

import (
	"context"
	"log"
	"sync"

	"github.com/nnurry/gopds/probabilistics/internal/config"

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

	log.Println("Successfully connected to Redis!")

})
