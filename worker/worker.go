package worker

import (
	"os"
	"os/signal"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var namespace = "ship_workers"
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "redis:6379")
	},
}

// Context ...
type Context struct {
	env *env.AppEnv
}

// Start ...
func Start(appEnv *env.AppEnv) error {
	url := os.Getenv("REDIS_URL")
	context := Context{env: appEnv}
	pool := work.NewWorkerPool(context, 10, namespace, &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", url)
		},
	})

	pool.Job(storeLogToAWS, (&context).StoreLogToAWS)

	pool.Start()
	defer pool.Stop()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	return nil
}
