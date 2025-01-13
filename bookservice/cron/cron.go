package cron

import (
	"fmt"
	"net/http"

	"github.com/robfig/cron/v3"
)

func SetupCronJobs() *cron.Cron {
	cronJob := cron.New()

	cronJob.AddFunc("0 0 * * *", func() {
		fmt.Println("Cron job start")

		resp, err := http.Get("http://orderservice:8080/orders/update-status")

		if err != nil {
			fmt.Println("Error making HTTP request:", err)
			return
		}

		defer resp.Body.Close()
	})

	return cronJob
}
