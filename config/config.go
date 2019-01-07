// Config package declares primitives to handle dynamic configuration.
package config

import (
	"github.com/go-ini/ini"
	"os"
)

// Config holds a complete set of dynamic configuration.
type Config struct {
	Environment string `json:"ENVIRONMENT"`
	AWSRegion   string `json:"AWSREGION"`
	//	Timeout       int64  `json:"TIMEOUT"`
	TelegramToken string `json:"TELEGRAMTOKEN"`
}

// GetConfigFromFile reads the configuration from an INI-style file and returns a Config struct.
// See tmb.conf.template for an example input file.
func GetConfigFromFile(configFile string) (*Config, error) {
	inicfg, err := ini.Load(configFile)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Environment:   inicfg.Section("").Key("ENVIRONMENT").String(),
		AWSRegion:     inicfg.Section("").Key("AWSREGION").String(),
		TelegramToken: inicfg.Section("").Key("TELEGRAMTOKEN").String(),
	}
	/*	cfg.Timeout, err = inicfg.Section("").Key("TIMEOUT").Int64()
		if err != nil {
			return nil, err
		}
	*/
	return &cfg, nil
}

// GetConfigFromENV reads the configuration from environment variables and returns a Config struct.
func GetConfigFromENV() (*Config, error) {
	config := Config{
		Environment:   os.Getenv("ENVIRONMENT"),
		AWSRegion:     os.Getenv("AWSREGION"),
		TelegramToken: os.Getenv("TELEGRAMTOKEN"),
	}

	/*	timeoutString := os.Getenv("TIMEOUT")
		if timeoutString == "" {
			timeoutString = "90"
		}
		timeout, err := strconv.ParseInt(timeoutString, 10, 64)
		if err != nil {
			return nil, err
		}
		config.Timeout = timeout
	*/return &config, nil
}
