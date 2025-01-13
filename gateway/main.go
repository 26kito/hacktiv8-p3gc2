package main

import "gateway/config"

// @title Book Library API
// version 1.0
// description API for book library
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	router := config.NewRouter()

	router.Logger.Fatal(router.Start(":8080"))
}
