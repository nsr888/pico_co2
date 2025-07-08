package main

import (
	"log"

	"pico_co2/internal/app"
)

func main() {
	log.SetFlags(0)

	application, err := app.New()
	if err != nil {
		log.Fatalf("application setup failed: %v", err)
	}

	application.Run()
}
