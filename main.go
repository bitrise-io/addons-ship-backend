package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/router"
	"github.com/bitrise-io/addons-ship-backend/worker"
	"github.com/bitrise-io/api-utils/logging"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	logger := logging.WithContext(nil)
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Failed to sync logger: %#v", err)
		}
	}()

	workerMode := os.Getenv("WORKER") == "true"
	if workerMode {
		appEnv, err := env.New(dataservices.GetDB())
		if err != nil {
			logger.Error("Failed to initialize Application Environment object for worker", zap.Any("error", err))
			os.Exit(1)
		}
		log.Println("Starting worker mode...")
		log.Fatal(errors.WithStack(worker.Start(appEnv)))
	} else {
		tracer.Start(tracer.WithServiceName("addons-ship"))
		defer tracer.Stop()

		err := dataservices.InitializeConnection(dataservices.ConnectionParams{}, true)
		if err != nil {
			logger.Error("Failed to initialize DB connection", zap.Any("error", err))
			os.Exit(1)
		}
		defer dataservices.Close()
		log.Println(" [OK] Database connection established")

		appEnv, err := env.New(dataservices.GetDB())
		if err != nil {
			logger.Error("Failed to initialize Application Environment object", zap.Any("error", err))
			os.Exit(1)
		}

		// Routing
		http.Handle("/", router.New(appEnv))

		log.Println("Starting - using port:", appEnv.Port)
		if err := http.ListenAndServe(":"+appEnv.Port, nil); err != nil {
			logger.Error("Failed to initialize Ship Addon Backend", zap.Any("error", err))
		}
	}
}
