package main

import (
	"log"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/router"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	appEnv := env.New(dataservices.GetDB())
	tracer.Start(tracer.WithServiceName("addons-ship"))
	defer tracer.Stop()

	err := dataservices.InitializeConnection(dataservices.ConnectionParams{}, true)
	if err != nil {
		appEnv.Logger.Error("Failed to initialize DB connection", zap.Any("error", err))
	}
	defer dataservices.Close()
	log.Println(" [OK] Database connection established")

	// Routing
	http.Handle("/", router.New(&appEnv))

	log.Println("Starting - using port:", appEnv.Port)
	if err := http.ListenAndServe(":"+appEnv.Port, nil); err != nil {
		appEnv.Logger.Error("Failed to initialize Ship Addon Backend", zap.Any("error", err))
	}
}
