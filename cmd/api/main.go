package main

import (
	"log"

	"github.com/6ivkin/test.git/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
