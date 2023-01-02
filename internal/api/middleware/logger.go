package middleware

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func JsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["request_path"] = params.Path
			log["request_method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006-01-02T15:04:05.00001Z")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}
