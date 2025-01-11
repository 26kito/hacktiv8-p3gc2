package main

import "gateway/config"

func main() {
	router := config.NewRouter()

	router.Logger.Fatal(router.Start(":8080"))
}
