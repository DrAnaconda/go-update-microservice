package main

import (
	"log"
	"update-microservice/internal/database/models"
	configuration_manager "update-microservice/packages/utils/configuration-manager"
	"update-microservice/packages/utils/database-utils"
)

func main() {
	_, err := configuration_manager.LoadConfiguration("config.json")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
		return
	}

	db, err := database_utils.OpenDatabaseConnection()
	if err != nil {
		log.Fatal("Failed to open database:", err)
		return
	}

	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
		return
	}
}
