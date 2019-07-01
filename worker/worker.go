package worker

import (
	"os"
	"os/signal"

	"github.com/bitrise-io/addons-ship-backend/env"
	redispkg "github.com/bitrise-io/addons-ship-backend/redis"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var namespace = "ship_workers"
var redisPool *redis.Pool

// Context ...
type Context struct {
	env *env.AppEnv
}

// Start ...
func Start(appEnv *env.AppEnv) error {
	urlStr := os.Getenv("REDIS_URL")
	context := Context{env: appEnv}
	redisPool = redispkg.NewPool(urlStr)
	pool := work.NewWorkerPool(context, 10, namespace, redisPool)

	pool.Job(storeLogToAWS, (&context).StoreLogToAWS)

	pool.Start()
	defer pool.Stop()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	return nil
}
