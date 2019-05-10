package main

import (
	"log"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/config"
	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/router"
	"github.com/bitrise-io/api-utils/logging"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	logger := logging.WithContext(nil)
	tracer.Start(tracer.WithServiceName("addons-ship"))
	defer tracer.Stop()

	err := dataservices.InitializeConnection(dataservices.ConnectionParams{})
	if err != nil {
		logger.Error("Failed to initialize DB connection", zap.Any("error", err))
	}
	defer dataservices.Close()
	log.Println(" [OK] Database connection established")

	// Routing
	http.Handle("/", router.New())

	appConfig := config.New(dataservices.GetDB())
	log.Println("Starting - using port:", appConfig.Port)
	if err := http.ListenAndServe(":"+appConfig.Port, nil); err != nil {
		logger.Error("Failed to initialize Ship Addon Backend", zap.Any("error", err))
	}
}
