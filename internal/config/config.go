package config

type Config struct {
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
	DBName         string
	LogFile        string
	RequestLimit   int
	MaxRLimitToken int
	JwtKey         string
	Address        string
	Port           string
}

var AppConfig Config
