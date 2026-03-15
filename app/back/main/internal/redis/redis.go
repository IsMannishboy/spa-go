package redis

import (
	"context"
	c "main/internal/config"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Cash struct {
	Rdb *redis.Client
}

func GetRedisConn(Server *c.Server) *Cash {
	rdb := redis.NewClient(&redis.Options{
		Addr: Server.Redis.Addr,

		DialTimeout:  time.Duration(Server.Redis.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(Server.Redis.DialTimeout) * time.Second,
		WriteTimeout: time.Duration(Server.Redis.DialTimeout) * time.Second,
		DB:           0,
	})
	var PingError error = nil
	for i := 0; i < 10; i++ {
		PingError = rdb.Ping(context.Background()).Err()
		if PingError == nil {
			break
		}
	}
	if PingError != nil {
		panic(PingError)
	}
	return &Cash{rdb}

}
