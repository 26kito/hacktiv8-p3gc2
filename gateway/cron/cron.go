package cron

import (
	"fmt"
	"net/http"

	"github.com/robfig/cron/v3"
)

func SetupCronJobs() *cron.Cron {
	cronJob := cron.New()

	// cronJob.AddFunc("0 0 * * *", func() {
	cronJob.AddFunc("* * * * *", func() {
		fmt.Println("Cron job start")

		resp, err := http.Get("http://localhost:8080/cron/book-update-status")

		if err != nil {
			fmt.Println("Error making HTTP request:", err)
			return
		}

		defer resp.Body.Close()
	})

	return cronJob
}
