package worker

import (
	"os"
	"os/signal"

	"github.com/bitrise-io/addons-ship-backend/env"
	redispkg "github.com/bitrise-io/addons-ship-backend/redis"
	"github.com/bitrise-io/api-utils/utils"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var namespace = "ship_workers"
var redisPool *redis.Pool

// Context ...
type Context struct {
	env *env.AppEnv
}

func init() {
	if redisPool == nil {
		redisPool = redispkg.NewPool(
			os.Getenv("REDIS_URL"),
			int(utils.GetInt64EnvWithDefault("WORKER_MAX_IDLE_CONNECTION", 50)),
			int(utils.GetInt64EnvWithDefault("WORKER_MAX_ACTIVE_CONNECTION", 1000)),
		)
	}
}

// Start ...
func Start(appEnv *env.AppEnv) error {
	context := Context{env: appEnv}
	pool := work.NewWorkerPool(context, 10, namespace, redisPool)

	pool.Job(storeLogToAWS, (&context).StoreLogToAWS)
	pool.Job(storeLogChunkToRedis, (&context).StoreLogChunkToRedis)
	pool.Job(copyUploadablesToNewAppVersion, (&context).CopyUploadablesToNewAppVersion)

	pool.Start()
	defer pool.Stop()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	return nil
}
