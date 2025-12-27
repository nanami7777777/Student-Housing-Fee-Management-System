package main

import (
	"dormsystem/config"
	"dormsystem/db"
	"dormsystem/handlers"
	"dormsystem/models"
	"dormsystem/router"
)

func main() {
	cfg := config.Load()
	db.Init(cfg.DBUrl,
		&models.ApartmentBuilding{},
		&models.DormRoom{},
		&models.Student{},
		&models.Payment{},
		&models.User{},
	)
	handlers.InitAuthData()
	r := router.SetupRouter()
	r.Run(cfg.HTTPPort)
}

