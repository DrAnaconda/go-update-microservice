package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"os"
	farmtracker "update-microservice/internal/farm-tracker"
	serverconfiguration "update-microservice/internal/server-configuration"
	configurationmanager "update-microservice/packages/utils/configuration-manager"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	serverconfiguration.ApplyGlobalRateLimit(e)

	config, err := configurationmanager.LoadConfiguration("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Load endpoint from DB (endpoint name >> target path/file name >> permissions)
	farmtracker.RegisterControllers(e)

	if os.Getenv("Debug") == "true" {
		e.Logger.Debug(e.StartTLS(config.ListenAddress, config.CertificatePath, config.CertificatePrivateKeyPath))
	} else {
		e.Logger.Info(e.StartTLS(config.ListenAddress, config.CertificatePath, config.CertificatePrivateKeyPath))
	}
}
