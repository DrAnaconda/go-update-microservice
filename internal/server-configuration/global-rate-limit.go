package server_configuration

import (
	"github.com/labstack/echo/v4"
	iplimiter "update-microservice/packages/middleware/ip-limiter"
)

func ApplyGlobalRateLimit(e *echo.Echo) {
	const requestPerSecond = 30

	e.Use(iplimiter.NewIPRateLimiterMiddleware(requestPerSecond))
}
