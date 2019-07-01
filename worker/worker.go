package worker

import (
	"os"
	"os/signal"

	"github.com/bitrise-io/addons-ship-backend/env"
	redispkg "github.com/bitrise-io/addons-ship-backend/redis"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
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
	urlStr := os.Getenv("REDIS_URL")
	context := Context{env: appEnv}
	pool := work.NewWorkerPool(context, 10, namespace, &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			url, err := redispkg.DialURL(urlStr)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			pass, err := redispkg.DialPassword(urlStr)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			c, err := redis.Dial("tcp", url, redis.DialPassword(pass))
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return c, nil
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
