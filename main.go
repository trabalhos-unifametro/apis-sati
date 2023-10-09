package main

import (
	"apis-sati/models"
	"apis-sati/routes"
	"apis-sati/utils"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func main() {
	app := routes.Routes()
	utils.LoadENVs()
	time.Local = time.UTC
	models.RunMigrations()

	if err := app.Listen(fmt.Sprint(":", viper.GetInt("PORT_APIS"))); err != nil {
		return
	}
}
