package ginext

import (
	"time"

	"github.com/gin-gonic/gin"

	"readygo/logger"
)

func GinLogger() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		// Process request
		c.Next()

		// Log only when path is not being skipped

		//Request := c.Request
		//isTerm := isTerm
		//Keys := c.Keys
		logFields := []logger.LogField{}
		// Stop timer
		logFields = append(logFields, logger.Time("start", start))
		TimeStamp := time.Now()
		logFields = append(logFields, logger.Duration("TimeStamp", TimeStamp.Sub(start)))
		logFields = append(logFields, logger.String("ClientIP", c.ClientIP()))
		logFields = append(logFields, logger.String("Method", c.Request.Method))
		logFields = append(logFields, logger.Int("StatusCode", c.Writer.Status()))
		logFields = append(logFields, logger.String("ErrorMessage", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		logFields = append(logFields, logger.Int("BodySize", c.Writer.Size()))

		if raw != "" {
			path = path + "?" + raw
		}
		logFields = append(logFields, logger.String("path", path))

		//记录日志
		logger.Info("[GIN]", logFields...)

	}
}
