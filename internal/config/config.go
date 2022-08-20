package config

import "github.com/spf13/viper"

type Config struct {
	Port            string `mapstructure:"PORT"`
	DBUrl           string `mapstructure:"DB_URL"`
	ApiKey          string `mapstructure:"API_KEY"`
	AuthSvcUrl      string `mapstructure:"AUTH_SVC_URL"`
	PassengerSvcUrl string `mapstructure:"PASSENGER_SVC_URL"`
	DriverSvcUrl    string `mapstructure:"DRIVER_SVC_URL"`
	GorushUrl       string `mapstructure:"GORUSH_URL"`
	FirebaseApiKey  string `mapstructure:"FIREBASE_API_KEY"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./internal/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
