package config

type Config struct {
	AppEnv      string
	AppPort     string
	GinMode     string
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   string
	CORSOrigin  string
	CORSMethods string
	CORSHeaders string
	ResendAPIKey string
	FromEmail   string
	FrontendURL string
}

var appConfig *Config

func GetConfig() *Config {
	if appConfig == nil {
		appConfig = &Config{
			AppEnv:      getEnv("APP_ENV", "development"),
			AppPort:     getEnv("APP_PORT", "8080"),
			GinMode:     getEnv("GIN_MODE", "debug"),
			JWTSecret:   getEnv("JWT_SECRET", "super-secret-jwt-key-here"),
			JWTExpiry:   getEnv("JWT_EXPIRY", "24h"),
			CORSOrigin:  getEnv("CORS_ALLOW_ORIGIN", "*"),
			CORSMethods: getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			CORSHeaders: getEnv("CORS_ALLOW_HEADERS", "*"),
			ResendAPIKey: getEnv("RESEND_API_KEY", ""),
			FromEmail:   getEnv("FROM_EMAIL", "onboarding@resend.dev"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		}
	}
	return appConfig
}