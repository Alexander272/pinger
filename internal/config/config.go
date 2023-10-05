package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Environment string       `yaml:"environment" env:"APP_ENV" env-default:"dev"`
		Pinger      PingerConfig `yaml:"pinger"`
		Bot         BotConfig    `yaml:"bot"`
		// Redis       RedisConfig
	}

	PingerConfig struct {
		Timeout   time.Duration `yaml:"timeout" env-default:"5s"`
		Interval  time.Duration `yaml:"interval" env-default:"1m"`
		Addresses []string      `yaml:"addresses"`
	}

	BotConfig struct {
		Server    string `env:"MOST_SERVER"`
		Token     string `env:"MOST_TOKEN"`
		ChannelId string `env:"MOST_CHANNEL_ID"`
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
