package farm_tracker

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
	"update-microservice/packages/dtos"
	basicsharedauth "update-microservice/packages/middleware/basic-shared-auth"
	configurationmanager "update-microservice/packages/utils/configuration-manager"
)

const baseControllerPath = "f-repo/farm-tracker/"

const metadataEndpoint = "metadata"
const updateFileEndpoint = "update"

func RegisterControllers(e *echo.Echo) {
	registerExplorerEndpoints(e)

	registerFarmTrackerEndpoints(e)
}

func registerFarmTrackerEndpoints(e *echo.Echo) {
	baseModulePath := fmt.Sprintf("%stracker/", baseControllerPath)

	e.GET(fmt.Sprintf("%s%s", baseModulePath, metadataEndpoint), func(c echo.Context) error {
		fileInfo, err := os.Stat(configurationmanager.Configuration.FarmTrackerUpdateFile)
		if err != nil {
			return err
		}

		modTime := fileInfo.ModTime()
		metadata := dtos.UpdateFileMetadata{
			Updated: modTime.Format(time.RFC3339),
		}

		return c.JSON(http.StatusOK, metadata)
	}, basicsharedauth.TokenAuthorizationMiddleware(configurationmanager.Configuration.FarmTrackerToken))

	e.GET(fmt.Sprintf("%s%s", baseModulePath, updateFileEndpoint), func(c echo.Context) error {
		return c.File(configurationmanager.Configuration.FarmTrackerUpdateFile)
	}, basicsharedauth.TokenAuthorizationMiddleware(configurationmanager.Configuration.FarmTrackerToken))
}

func registerExplorerEndpoints(e *echo.Echo) {
	baseModulePath := fmt.Sprintf("%sexplorer/", baseControllerPath)

	e.GET(fmt.Sprintf("%s%s", baseModulePath, metadataEndpoint), func(c echo.Context) error {
		fileInfo, err := os.Stat(configurationmanager.Configuration.FarmExplorerUpdateFile)
		if err != nil {
			return err
		}

		modTime := fileInfo.ModTime()
		metadata := dtos.UpdateFileMetadata{
			Updated: modTime.Format(time.RFC3339),
		}

		return c.JSON(http.StatusOK, metadata)
	}, basicsharedauth.TokenAuthorizationMiddleware(configurationmanager.Configuration.FarmTrackerToken))

	e.GET(fmt.Sprintf("%s%s", baseModulePath, updateFileEndpoint), func(c echo.Context) error {
		return c.File(configurationmanager.Configuration.FarmExplorerUpdateFile)
	}, basicsharedauth.TokenAuthorizationMiddleware(configurationmanager.Configuration.FarmTrackerToken))
}
