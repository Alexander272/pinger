package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Environment string       `yaml:"environment" env:"APP_ENV" env-default:"dev"`
		LogLevel    string       `yaml:"log_level" env-default:"info"`
		LogSource   bool         `yaml:"log_source" env-default:"false"`
		Http        HttpConfig   `yaml:"http"`
		Pinger      PingerConfig `yaml:"pinger"`
		Bot         BotConfig    `yaml:"bot"`
		Postgres    PostgresConfig
		Redis       RedisConfig
	}

	HttpConfig struct {
		Host           string        `yaml:"host" env:"HOST" env-default:"localhost"`
		Port           string        `yaml:"port" env:"PORT" env-default:"8080"`
		ReadTimeout    time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"10s"`
		WriteTimeout   time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"10s"`
		MaxHeaderBytes int           `yaml:"max_header_bytes" env-default:"1"`
	}

	PingerConfig struct {
		Count     int                `yaml:"count" env-default:"5"`
		Interval  time.Duration      `yaml:"interval" env-default:"0.1s"`
		Timeout   time.Duration      `yaml:"timeout" env-default:"1s"`
		IP        string             `yaml:"ip" env:"IP"`
		Rtt       time.Duration      `yaml:"rtt" env-default:"50ms"`
		Addresses []*AddressesConfig `yaml:"addresses"`
	}

	AddressesConfig struct {
		Interval time.Duration `yaml:"interval"`
		List     []*Address    `yaml:"list"`
	}
	Address struct {
		Ip   string        `yaml:"ip"`
		Name string        `yaml:"name"`
		Rtt  time.Duration `yaml:"rtt"`
	}

	BotConfig struct {
		Server    string `env:"MOST_SERVER"`
		Token     string `env:"MOST_TOKEN"`
		ChannelId string `env:"MOST_CHANNEL_ID" yaml:"channel_id"`
	}

	PostgresConfig struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST"`
		Port     string `yaml:"port" env:"POSTGRES_PORT"`
		Username string `yaml:"username" env:"POSTGRES_NAME"`
		Password string `env:"POSTGRES_PASSWORD"`
		DbName   string `yaml:"db_name" env:"POSTGRES_DB"`
		SSLMode  string `yaml:"ssl_mode" env:"POSTGRES_SSL"`
	}

	RedisConfig struct {
		Host     string `yaml:"host" env:"REDIS_HOST"`
		Port     string `yaml:"port" env:"REDIS_PORT"`
		DB       int    `yaml:"db" env:"REDIS_DB"`
		Password string `env:"REDIS_PASSWORD"`
	}
)

func Init(path string) (*Config, error) {
	var conf Config

	if err := cleanenv.ReadConfig(path, &conf); err != nil {
		return nil, fmt.Errorf("failed to read config file. error: %w", err)
	}

	// if err := cleanenv.ReadEnv(&conf); err != nil {
	// 	return nil, fmt.Errorf("failed to read env file. error: %w", err)
	// }

	return &conf, nil
}

// TODO дописать получение новых данных для конфига
// func (c *Config) Update() error {
// 	out, err := yaml.Marshal(&c)
// 	if err != nil {
// 		// logger.Errorf("failed marshal yaml. error: %s", err.Error())
// 		return fmt.Errorf("failed marshal yaml. error: %w", err)
// 	}

// 	err = os.WriteFile("modified.yaml", out, 0777)
// 	if err != nil {
// 		logger.Errorf("Problem updating file: %s", err)
// 		return fmt.Errorf("problem updating file: %w", err)
// 	}

// 	return nil
// }
