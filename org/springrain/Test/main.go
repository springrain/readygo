package main

import (
	"goshop/org/springrain/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(GinLogger())
	//r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

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
		logger.Info("ginlog", logFields...)

	}
}
