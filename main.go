package main

import (
	"dumbky/internal/app"
	"dumbky/internal/log"
)

func main() {
	log.Init()
	app.Run()
}

