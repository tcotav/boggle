package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// TODO set the global log level to info -- probably better sent in via config than h/coded
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// RequestLogger is a gin middleware that logs requests in json format and prints out pretty standard http access logs
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		client_ip := c.ClientIP()
		user_agent := c.Request.UserAgent()
		method := c.Request.Method
		path := c.Request.URL.Path

		t := time.Now()
		c.Next() // this hands off to the next handler in the chain
		latency := float32(time.Since(t).Seconds())
		status := c.Writer.Status()

		log.Info().Str("client_ip", client_ip).Str("user_agent", user_agent).Str("method", method).Str("path", path).Float32("latency", latency).Int("status", status).Msg("")
	}
}
