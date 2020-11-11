package middleware

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// author zhasulan
// created on 30.09.20 20:31

// Copied from github.com/toorop/gin-logrus as modified for FCB

// Logger is the logrus logger handler
func Logger(logger logrus.FieldLogger, notLogged ...string) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, p := range notLogged {
			skip[p] = struct{}{}
		}
	}

	return func(context *gin.Context) {
		// other handler can change context.Path so:
		path := context.Request.URL.Path
		start := time.Now()
		context.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := context.Writer.Status()
		clientIP := context.ClientIP()
		clientUserAgent := context.Request.UserAgent()
		referer := context.Request.Referer()
		dataLength := context.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		if _, ok := skip[path]; ok {
			return
		}

		entry := logger.WithFields(logrus.Fields{
			//"timestamp":  time.Now().Format(TimestampFormat),
			"addr": clientIP,
			// Other fields not required FCB Kubernetes
			//"hostname":   hostname,
			//"statusCode": statusCode,
			//"latency":    latency, // time to process
			//"method":     context.Request.Method,
			//"path":       path,
			//"referer":    referer,
			//"dataLength": dataLength,
			//"userAgent":  clientUserAgent,
		})

		if len(context.Errors) > 0 {
			entry.Error(context.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", // %s - %s [%s]
				//clientIP,
				hostname,
				//time.Now().Format(TimestampFormat),
				context.Request.Method,
				path,
				statusCode,
				dataLength,
				referer,
				clientUserAgent,
				latency,
			)
			if statusCode >= http.StatusInternalServerError {
				entry.Error(msg)
			} else if statusCode >= http.StatusBadRequest {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
