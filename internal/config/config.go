package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		URL string
	}
	Auth struct {
		JWTSecret     string `mapstructure:"jwt_secret"`
		HMACSecret    string `mapstructure:"hmac_secret"`
		EncryptionKey string `mapstructure:"encryption_key"`
	}
	SMTP struct {
		Host string
		Port int
		User string
		Pass string
		From string
	}
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
