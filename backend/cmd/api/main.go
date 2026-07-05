package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Varfa/GarageHub/internal/config"
	"github.com/Varfa/GarageHub/internal/database"
	"github.com/Varfa/GarageHub/internal/handler"
	"github.com/Varfa/GarageHub/internal/repository"
	"github.com/Varfa/GarageHub/internal/router"
	"github.com/Varfa/GarageHub/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	dbPool, err := database.ConnectPostgres(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	clientRepo := repository.NewClientRepository(dbPool)
	clientService := service.NewClientService(clientRepo)
	clientHandler := handler.NewClientHandler(clientService)

	mux := router.SetupRoutes(clientHandler)
	addr := ":" + cfg.Server.Port

	fmt.Println("Server started on", addr)

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}

}
