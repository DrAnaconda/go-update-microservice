package ip_limiter

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
	"time"
)

type IPRateLimiter struct {
	clientLimiters map[string]*rate.Limiter
	limitDuration  time.Duration
}

func NewIPRateLimiter(seconds int) *IPRateLimiter {
	return &IPRateLimiter{
		clientLimiters: make(map[string]*rate.Limiter),
		limitDuration:  time.Duration(seconds) * time.Second,
	}
}

func NewIPRateLimiterMiddleware(seconds int) echo.MiddlewareFunc {
	rl := &IPRateLimiter{
		clientLimiters: make(map[string]*rate.Limiter),
		limitDuration:  time.Duration(seconds) * time.Second,
	}

	return rl.Limit
}

func (rl *IPRateLimiter) Limit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		clientIP := c.RealIP()

		if _, clientExists := rl.clientLimiters[clientIP]; !clientExists {
			// Allow 1 request per the specified duration
			rl.clientLimiters[clientIP] = rate.NewLimiter(rate.Every(rl.limitDuration), 20)
		}

		if rl.clientLimiters[clientIP].Allow() == false {
			return echo.ErrTooManyRequests
		}

		return next(c)
	}
}
