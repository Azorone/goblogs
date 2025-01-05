package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"mastcat/internal/model"
	"mastcat/pkg/util"
	"net/http"
	"time"
)

func CrawlerFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		if userAgent == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User-Agent header required"})
			return
		}
		// List of known crawler User-Agent substrings
		//crawlers := []string{"Googlebot", "Bingbot", "Slurp", "DuckDuckBot", "Baiduspider", "YandexBot", "Sogou"}
		//startTime := time.Now()
		c.Next()
	}
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return util.JwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		accessLog := model.AccessLog{
			AccessTime:   startTime,
			IP:           c.ClientIP(),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			ResponseTime: time.Duration(time.Since(startTime).Milliseconds()),
		}
		fl := false
		if c.Writer.Status() >= 400 {
			fl = true
			log.Printf("Error: Status %d, Path %s", c.Writer.Status(), c.Request.URL.Path)
			if c.Writer.Status() == http.StatusNotFound {
				c.JSON(http.StatusOK, util.NewResponse(404, "Not Found", nil))
				accessLog.ErrorMsg = "Not Found"
			} else if c.Writer.Status() == http.StatusTooManyRequests {
				c.JSON(http.StatusOK, util.NewResponse(429, "Too Many Requests", nil))
				accessLog.ErrorMsg = "Too Many Requests"
			} else if c.Writer.Status() == http.StatusUnauthorized {
				c.JSON(http.StatusOK, util.NewResponse(401, "Unauthorized", nil))
				accessLog.ErrorMsg = "Unauthorized"
			} else if c.Writer.Status() == http.StatusForbidden {
				c.JSON(http.StatusOK, util.NewResponse(403, "Forbidden", nil))
				accessLog.ErrorMsg = "Forbidden"
			} else {
				c.JSON(http.StatusOK, util.NewResponse(500, "Server Error", nil))
				accessLog.ErrorMsg = "Server Error"
			}
		}
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("Error: %v", e.Err)
				accessLog.ErrorMsg = accessLog.ErrorMsg + "|" + e.Err.Error()
			}
			if !fl {

				c.JSON(http.StatusOK, util.NewResponse(500, "Server Error", nil))
			}

		}
		util.AppendLogBuff(accessLog)
	}
}
func AccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the current time
		startTime := time.Now()

		// Process request
		c.Next()

		// Create an access log entry
		accessLog := model.AccessLog{
			AccessTime:   startTime,
			IP:           c.ClientIP(),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			ResponseTime: time.Duration(time.Since(startTime).Milliseconds()),
		}
		util.AppendLogBuff(accessLog)
	}
}
func LimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.ClientIP()
		if !util.GetLimiter(c.ClientIP()).Allow() {
			startTime := time.Now()
			// Create an access log entry
			accessLog := model.AccessLog{
				AccessTime:   startTime,
				IP:           c.ClientIP(),
				Method:       c.Request.Method,
				Path:         c.Request.URL.Path,
				ResponseTime: time.Duration(time.Since(startTime).Milliseconds()),
				ErrorMsg:     "too many requests",
			}
			util.AppendLogBuff(accessLog)
			if err := c.AbortWithError(http.StatusTooManyRequests, fmt.Errorf("too many requests")); err != nil {
				log.Printf("Error: %v", err)
			}

		}
		c.Next()
	}
}
