package redis

import (
	c "auth/internal/config"
	"time"

	redis "github.com/redis/go-redis/v9"
)

func GetRedisConn(Server *c.Server) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         Server.Redis.Host,
		DialTimeout:  time.Duration(Server.Redis.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(Server.Redis.DialTimeout) * time.Second,
		WriteTimeout: time.Duration(Server.Redis.DialTimeout) * time.Second,
		DB:           0,
	})
	var PingError error
	for i := 0; i < 10; i++ {
		PingError = rdb.Ping(nil).Err()
		if PingError == nil {
			break
		}
	}
	if PingError != nil {
		panic(PingError)
	}
	return rdb

}
