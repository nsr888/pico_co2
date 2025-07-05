package main

import (
	"log"

	"pico_co2/internal/app"
)

func main() {
	log.SetFlags(0)

	config := app.Config{
		IsAdvancedSetup: false,
	}

	application, err := app.New(config)
	if err != nil {
		log.Fatalf("application setup failed: %v", err)
	}

	application.Run()
}
