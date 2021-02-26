package mid

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// DbConfig ....
type DbConfig struct {
	Name     string `json:"database.host"`
	Port     string `json:"database.port"`
	User     string `json:"database.user"`
	Password string `json:"database.pass"`
	Dbname   string `json:"database.dbname"`
}

// LoadConfig ....
func LoadConfig(path string) (config DbConfig, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()

	if err != nil {
		d := err.Error()
		fmt.Println(d)
		return DbConfig{}, errors.WithStack(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return DbConfig{}, errors.WithStack(err)
	}
	return DbConfig{}, nil
}
