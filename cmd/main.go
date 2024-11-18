package main

import (
	"github.com/tittuvarghese/ss-go-core/config"
	"github.com/tittuvarghese/ss-go-core/logger"
	"github.com/tittuvarghese/ss-go-order-management-service/constants"
	"github.com/tittuvarghese/ss-go-order-management-service/core/database"
	"github.com/tittuvarghese/ss-go-order-management-service/core/handler"
	"github.com/tittuvarghese/ss-go-order-management-service/models"
)

func main() {
	log := logger.NewLogger(constants.ModuleName)
	log.Info("Initialising Customer Service Module")

	// Config Management
	configManager := config.NewConfigManager(config.DEFAULT_CONFIG_PATH)
	configManager.Enable()

	if configManager.GetBool(constants.OtelEnableEnv) {
		serviceName := configManager.GetString(constants.OtelServiceNameEnv)
		collectorUrl := configManager.GetString(constants.OtelCollectorEnv)
		insecureMode := configManager.GetBool(constants.OtelInsecureModeEnv)
		otel.NewTraceProvider(serviceName, collectorUrl, insecureMode)
	}

	// DB Handling
	dbConn := configManager.GetString(constants.DatabaseUrlEnvName)

	log.Info("DB Connection String " + dbConn)

	dbInstance, err := database.NewRelationalDatabase(dbConn)
	if err != nil {
		log.Error("Error initialising relational db", err)
	}

	err = dbInstance.Instance.Open()
	if err != nil {
		log.Error("Error opening relational db", err)
	}

	err = dbInstance.Instance.AutoMigrate(models.Order{}, models.Item{}, models.Address{})
	if err != nil {
		log.Error("Error performing auto migration for db", err)
	}

	server := handler.NewGrpcServer()
	server.RdbInstance = dbInstance
	server.Run(constants.GrpcServerPort)
}
