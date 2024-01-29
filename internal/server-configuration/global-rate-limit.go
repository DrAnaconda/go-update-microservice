package server_configuration

import (
	"github.com/labstack/echo/v4"
	iplimiter "update-microservice/packages/middleware/ip-limiter"
	configuration_manager "update-microservice/packages/utils/configuration-manager"
)

func ApplyGlobalRateLimit(e *echo.Echo) {
	e.Use(iplimiter.NewIPRateLimiterMiddleware(configuration_manager.Configuration.GlobalRequestsPerSecond))
}
