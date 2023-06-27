package config

import (
	"os"
)

type dbConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

type jwtConfig struct {
	ExpTimeMinutes string
	SecretString   string
	JWTIssuer      string
}

type envConfig struct {
	Mode string
}

type cloudinaryConfig struct {
	CloudName    string
	APIKey       string
	APISecret    string
	UploadFolder string
}

type AppConfig struct {
	DBConfig         dbConfig
	JWTConfig        jwtConfig
	ENVConfig        envConfig
	CloudinaryConfig cloudinaryConfig
}

func getEnv(key, defaultVal string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultVal
	}
	return env
}

func InitConfig() AppConfig {
	config := AppConfig{
		DBConfig: dbConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "dummy_db"),
			Port:     getEnv("DB_PORT", "5432"),
		},

		JWTConfig: jwtConfig{
			ExpTimeMinutes: getEnv("JWT_EXPIRATION", "15"),
			SecretString:   getEnv("SECRET_KEY", "very-secret-key"),
			JWTIssuer:      getEnv("JWT_ISSUER", "localhost"),
		},

		ENVConfig: envConfig{
			Mode: getEnv("ENV_MODE", "testing"),
		},

		CloudinaryConfig: cloudinaryConfig{
			CloudName:    getEnv("CLOUDINARY_CLOUD_NAME", "dxfq3iotg"),
			APIKey:       getEnv("CLOUDINARY_API_KEY", "111 111 111 111"),
			APISecret:    getEnv("CLOUDINARY_API_SECRET", "111 111 111 111"),
			UploadFolder: getEnv("CLOUDINARY_UPLOAD_FOLDER", "final-project"),
		},
	}
	return config
}
