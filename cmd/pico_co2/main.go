package main

import (
	"pico_co2/internal/app"
)

func main() {
	cfg := app.DefaultConfig()
	application, err := app.New(cfg)
	if err != nil {
		println("error creating application:", err)
		return
	}

	application.Run()
}
