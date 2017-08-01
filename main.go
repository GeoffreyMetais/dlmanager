package main

import (
	"github.com/GeoffreyMetais/dlmanager/api"
	"github.com/GeoffreyMetais/dlmanager/db"
)

func main() {
	defer db.PrepareDb().Close()
	api.Run()
}
