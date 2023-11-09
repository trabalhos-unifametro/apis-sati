package main

import (
	"apis-sati/models"
	"apis-sati/routes"
	"apis-sati/utils"
	"fmt"
	"time"
)

func init() {
	utils.LoadENVs()
}

func main() {
	app := routes.Routes()
	time.Local = time.UTC
	models.RunMigrations()

	if err := app.Listen(fmt.Sprint(":", utils.GodotEnv("PORT_APIS"))); err != nil {
		return
	}
}
