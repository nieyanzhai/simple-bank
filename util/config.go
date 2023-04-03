package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (Config, error) {
	var c Config
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return c, err
	}
	err = viper.Unmarshal(&c)
	return c, err
}
