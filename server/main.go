package main

import (
	"log"

	"dormsystem/config"
	"dormsystem/db"
	"dormsystem/handlers"
	"dormsystem/models"
	"dormsystem/mq"
	"dormsystem/router"
)

func main() {
	cfg := config.Load()
	if err := mq.Init(cfg.MQUrl); err != nil {
		log.Println("mq init error", err)
	}
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
