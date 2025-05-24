package main

import (
	"dumbky/internal/app"
	"dumbky/internal/db"
	"dumbky/internal/log"
)

func main() {
	log.Init()
	db.Init()
	app.Run()
}
