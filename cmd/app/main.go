package main

import (
	"os"

	"github.com/cloudingcity/todo/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
