package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"mastcat/internal/config"
	"mastcat/internal/database"
	"mastcat/internal/model"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex
var JwtKey = []byte(config.AppConfig.JwtKey)
var LogBuffer []model.AccessLog
var LogBufferMu sync.Mutex

// StructToMap converts a struct to a map, ignoring empty fields
func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			continue
		}
		if (field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface) && field.IsNil() {
			continue
		}
		result[fieldType.Name] = field.Interface()
	}
	return result
}

// NewResponse creates a new response object
func NewResponse(code int, msg string, data interface{}) Response {
	return Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// GenerateToken generates a JWT token for a user
func GenerateToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a JWT token and returns the username
func ParseToken(tokenString string) (string, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	username := (*claims)["username"].(string)
	return username, nil
}

// IsValidEmail checks if an email address is valid
func IsValidEmail(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return strings.Contains(email, "@")
}

// 每秒允许5次请求，最多存储10个令牌
func GetLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(
			rate.Limit(config.AppConfig.RequestLimit),
			config.AppConfig.MaxRLimitToken)
		visitors[ip] = limiter
	}

	return limiter
}
func CleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, limiter := range visitors {
			if limiter.Allow() { // 如果1分钟内没有请求到达
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func SaveLogs(logs []model.AccessLog) {
	for _, log := range logs {
		if err := database.DB.Create(&log).Error; err != nil {

			continue
		}
	}
}
func GetGinLog(c *gin.Context) {

}
func AppendLogBuff(log model.AccessLog) {
	LogBufferMu.Lock()
	LogBuffer = append(LogBuffer, log)
	if len(LogBuffer) >= 40 {
		go FlushLogBuffer()
	}
	defer LogBufferMu.Unlock()
}

func FlushLogBuffer() {

	LogBufferMu.Lock()
	defer LogBufferMu.Unlock()

	if len(LogBuffer) > 0 {
		SaveLogs(LogBuffer)
		LogBuffer = LogBuffer[:0] // Clear the buffer
	}
}
func init() {
	go func() {
		for {
			time.Sleep(2 * time.Minute)
			FlushLogBuffer()
		}
	}()

	go CleanupVisitors()
}
