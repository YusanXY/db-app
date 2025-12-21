package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	File     FileConfig     `mapstructure:"file"`
	App      AppConfig      `mapstructure:"app"`
}

type ServerConfig struct {
	Port         string `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int    `mapstructure:"expires_in"`
	RefreshIn int    `mapstructure:"refresh_in"`
}

type FileConfig struct {
	UploadPath  string   `mapstructure:"upload_path"`
	MaxSize     int64    `mapstructure:"max_size"`
	AllowedExt  []string `mapstructure:"allowed_ext"`
}

type AppConfig struct {
	Name  string `mapstructure:"name"`
	Env   string `mapstructure:"env"`
	Debug bool   `mapstructure:"debug"`
}

var GlobalConfig *Config

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	// 环境变量支持
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// 从环境变量读取配置
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "SERVER_MODE")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.db_name", "DB_NAME")
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.expires_in", "JWT_EXPIRES_IN")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 设置默认值
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.Mode == "" {
		config.Server.Mode = "debug"
	}
	if config.JWT.Secret == "" {
		config.JWT.Secret = "default-secret-key-change-in-production"
	}

	GlobalConfig = &config
	return &config, nil
}

func GetConfig() *Config {
	if GlobalConfig == nil {
		// 尝试从环境变量加载
		config, err := LoadConfig()
		if err != nil {
			// 使用默认配置
			config = &Config{
				Server: ServerConfig{
					Port: getEnv("SERVER_PORT", "8080"),
					Mode: getEnv("SERVER_MODE", "debug"),
				},
				Database: DatabaseConfig{
					Host:     getEnv("DB_HOST", "localhost"),
					Port:     5432,
					User:     getEnv("DB_USER", "dbapp"),
					Password: getEnv("DB_PASSWORD", "dbapp123"),
					DBName:   getEnv("DB_NAME", "dbapp"),
				},
				Redis: RedisConfig{
					Host: getEnv("REDIS_HOST", "localhost"),
					Port: 6379,
				},
				JWT: JWTConfig{
					Secret:    getEnv("JWT_SECRET", "default-secret-key"),
					ExpiresIn: 3600,
				},
			}
		}
		GlobalConfig = config
	}
	return GlobalConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

