package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"os"
	"update-microservice/internal/controllers"
	serverconfiguration "update-microservice/internal/server-configuration"
	configurationmanager "update-microservice/packages/utils/configuration-manager"
)

// TODO: License handling for each product
// TODO: Config from DB OR Local file
// TODO: Rate limiting options from DB for each product

func main() {
	config, err := configurationmanager.LoadConfiguration("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	serverconfiguration.ApplyGlobalRateLimit(e)

	controllers.RegisterUpdateEndpoints(e)

	if os.Getenv("Debug") == "true" {
		e.Logger.Debug(e.StartTLS(config.ListenAddress, config.CertificatePath, config.CertificatePrivateKeyPath))
	} else {
		e.Logger.Info(e.StartTLS(config.ListenAddress, config.CertificatePath, config.CertificatePrivateKeyPath))
	}
}
