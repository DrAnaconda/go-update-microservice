package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
	"update-microservice/internal/database/models"
	"update-microservice/internal/database/repos"
	"update-microservice/packages/dtos"
	basicsharedauth "update-microservice/packages/middleware/basic-shared-auth"
	ip_limiter "update-microservice/packages/middleware/ip-limiter"
	database_utils "update-microservice/packages/utils/database-utils"
)

const baseControllerPath = "f-repo/farm-tracker"

const metadataEndpoint = "metadata"
const updateFileEndpoint = "update"

var products []models.Product

func RegisterUpdateEndpoints(e *echo.Echo) {
	db, err := database_utils.OpenDatabaseConnection()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to open db connection: %w", err))
		return
	}

	productRepo := repos.NewProductRepository(db)

	products, err = productRepo.GetAllProducts()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to load products: %w", err))
		return
	}

	createMetadataEndpoints(e)
	createUpdateFileEndpoints(e)
}

func createUpdateFileEndpoints(e *echo.Echo) {
	for _, product := range products {
		baseModulePath := fmt.Sprintf("%s/%s/", baseControllerPath, product.Name)

		updateRateLimiter := ip_limiter.NewIPRateLimiterMiddleware(int((5 * time.Minute).Seconds()))

		routeGroup := e.Group(fmt.Sprintf("%s%s", baseModulePath, updateFileEndpoint))
		routeGroup.Use(updateRateLimiter)
		routeGroup.Use(basicsharedauth.TokenAuthorizationMiddleware(product.Password))

		finalPath := fmt.Sprintf("%s%s", baseModulePath, updateFileEndpoint)

		e.GET(finalPath, func(c echo.Context) error {
			return c.File(product.PathToProduct)
		})
	}
}

func createMetadataEndpoints(e *echo.Echo) {
	for _, product := range products {
		baseModulePath := fmt.Sprintf("%s/%s/", baseControllerPath, product.Name)

		e.GET(fmt.Sprintf("%s%s", baseModulePath, metadataEndpoint), func(c echo.Context) error {
			fileInfo, err := os.Stat(product.PathToProduct)
			if err != nil {
				return err
			}

			modTime := fileInfo.ModTime()
			metadata := dtos.UpdateFileMetadata{
				Updated: modTime.Format(time.RFC3339),
			}

			return c.JSON(http.StatusOK, metadata)
		}, basicsharedauth.TokenAuthorizationMiddleware(product.Password))
	}
}
