package main

import (
	"gateway/config"
	"gateway/cron"
)

// @title Book Library API
// version 1.0
// description API for book library
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	router := config.NewRouter()

	cronJob := cron.SetupCronJobs()
	cronJob.Start()

	router.Logger.Fatal(router.Start(":8080"))
}
