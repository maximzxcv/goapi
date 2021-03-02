package mid

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// DbConfig ....
// type DbConfig struct {
// 	Name     string
// 	Port     string
// 	User     string
// 	Password string
// 	Dbname   string
// }

// AppConfig .....
type AppConfig struct {
	Db DbConfig
}

// DbConfig ....
type DbConfig struct {
	Host     string `json:"database.host"`
	Port     string `json:"database.port"`
	User     string `json:"database.user"`
	Password string `json:"database.pass"`
	Dbname   string `json:"database.dbname"`
}

// NewTestConfig uses for testing
func NewTestConfig() *DbConfig {
	return &DbConfig{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "goapitestpass",
		Dbname:   "postgres",
	}
}

// ConnectinString ....
func (dbConfig *DbConfig) ConnectinString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Dbname)
}

// LoadConfig ....
func LoadConfig(path string) error { //} (*AppConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		//return nil, errors.Wrap(err, "Error reading configuration")
		return errors.Wrap(err, "Error reading configuration")
	}

	var appConfig AppConfig

	return viper.Unmarshal(&appConfig)
}

// GetDbConfig ....
func GetDbConfig() (*DbConfig, error) {

	conf := func(key string) string {
		return viper.GetString(`database.` + key)
	}

	return &DbConfig{
		Host:     conf("host"),
		Port:     conf("port"),
		User:     conf("user"),
		Password: conf("pass"),
		Dbname:   conf("dbname"),
	}, nil
}
