package app

import (
	"fmt"

	"github.com/spf13/viper"
)

type App struct {
	// Db *gorm.DB
}

var MyName *string

func SetName() {
	// a := "DDD"
	*MyName = "Sdfds"
}

func (app App) Version() string {

	return "0.0.1"
}

func (a App) Start() {
	fmt.Println("App Started")
}

func getAndSetConfig() {
	viper.SetConfigName("default") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func New() App {
	getAndSetConfig()
	return App{}
}
