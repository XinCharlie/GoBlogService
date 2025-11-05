package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Logging  LoggingConfig  `yaml:"logging"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	SSLMode      string `yaml:"sslmode"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type JWTConfig struct {
	Secret                string        `yaml:"key"`
	Issuer                string        `yaml:"issuer"`
	Audience              string        `yaml:"audience"`
	ExpirationHours       int           `yaml:"accessTokenExpiryHours"`
	RefreshExpirationDays int           `yaml:"refreshTokenExpiryDays"`
	Expiration            time.Duration `yaml:"-"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type KafkaConfig struct {
	KafkaBrokers  string `yaml:"kafka_brokers"`
	KafkaTopic    string `yaml:"kafka_topic"`
	KafkaGroupID  string `yaml:"kafka_groupID"`
	KafkaUsername string `yaml:"kafka_username"`
	KafkaPassword string `yaml:"kafka_password"`
}

func Load(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	// 计算过期时间
	config.JWT.Expiration = time.Duration(config.JWT.ExpirationHours) * time.Hour
	return config, nil
}
