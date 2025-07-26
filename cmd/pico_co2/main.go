package main

import (
	"pico_co2/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		println("error creating application:", err)
	}

	application.Run()
}
