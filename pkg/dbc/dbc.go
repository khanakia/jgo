package dbc

import (
	"fmt"

	loggerz "github.com/khanakia/jgo/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres" //Gorm postgres dialect interface
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbc struct {
	Db *gorm.DB
}

type Config struct {
	Logger loggerz.Logger
}

func New(config Config) Dbc {
	username := viper.GetString("database.user")
	password := viper.GetString("database.password")
	databaseName := viper.GetString("database.name")
	databaseHost := viper.GetString("database.host")
	databasePort := viper.GetString("database.port")

	//Define DB connection string
	dbDSN := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", databaseHost, username, databaseName, password, databasePort)

	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		// fmt.Println("error", err)
		config.Logger.SugarLogger.Error(err)
		panic(err)
	}
	return Dbc{
		Db: db,
	}
}
