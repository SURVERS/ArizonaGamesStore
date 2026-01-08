package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
	blocked  map[string]time.Time
	limit    int
	window   time.Duration
	blockDuration time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		blocked:  make(map[string]time.Time),
		limit:    limit,
		window:   window,
		blockDuration: 15 * time.Minute,
	}

	go rl.cleanup()

	return rl
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()

		for ip, times := range rl.requests {
			var valid []time.Time
			for _, t := range times {
				if now.Sub(t) < rl.window {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = valid
			}
		}

		for ip, blockedUntil := range rl.blocked {
			if now.After(blockedUntil) {
				delete(rl.blocked, ip)
			}
		}

		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if blockedUntil, exists := rl.blocked[ip]; exists {
		if time.Now().Before(blockedUntil) {
			return false
		}
		delete(rl.blocked, ip)
	}

	now := time.Now()
	var valid []time.Time

	if times, exists := rl.requests[ip]; exists {
		for _, t := range times {
			if now.Sub(t) < rl.window {
				valid = append(valid, t)
			}
		}
	}

	if len(valid) >= rl.limit {
		rl.blocked[ip] = now.Add(rl.blockDuration)
		return false
	}

	valid = append(valid, now)
	rl.requests[ip] = valid

	return true
}

func (rl *rateLimiter) isBlocked(ip string) (bool, time.Duration) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if blockedUntil, exists := rl.blocked[ip]; exists {
		if time.Now().Before(blockedUntil) {
			remaining := time.Until(blockedUntil)
			return true, remaining
		}
	}
	return false, 0
}

var (
	registerLimiter = newRateLimiter(3, 1*time.Hour)
	loginLimiter    = newRateLimiter(5, 5*time.Minute)
	verifyLimiter   = newRateLimiter(10, 10*time.Minute)
)

func RateLimitRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if blocked, remaining := registerLimiter.isBlocked(ip); blocked {
			minutes := int(remaining.Minutes())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Слишком много попыток регистрации. Ваш IP заблокирован на %d мин.", minutes),
			})
			c.Abort()
			return
		}

		if !registerLimiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Слишком много попыток регистрации. Попробуйте позже (максимум 3 регистрации в час)",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RateLimitLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if blocked, remaining := loginLimiter.isBlocked(ip); blocked {
			minutes := int(remaining.Minutes())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Слишком много неудачных попыток входа. Ваш IP заблокирован на %d мин.", minutes),
			})
			c.Abort()
			return
		}

		if !loginLimiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Слишком много попыток входа. Попробуйте через несколько минут (максимум 5 попыток за 5 минут)",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RateLimitVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if blocked, remaining := verifyLimiter.isBlocked(ip); blocked {
			minutes := int(remaining.Minutes())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Слишком много попыток верификации. Ваш IP заблокирован на %d мин.", minutes),
			})
			c.Abort()
			return
		}

		if !verifyLimiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Слишком много попыток верификации. Попробуйте через несколько минут",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
