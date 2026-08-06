package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Varfa/GarageHub/internal/config"
	"github.com/Varfa/GarageHub/internal/database"
	"github.com/Varfa/GarageHub/internal/handler"
	"github.com/Varfa/GarageHub/internal/i18n"
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

	translator, err := i18n.NewManager(
		[]string{"en", "lt", "uk", "ru"},
		"en",
	)
	if err != nil {
		log.Fatal(err)
	}

	handler.SetTranslator(translator)

	loginHandler := handler.NewLoginHandler(translator)

	clientRepository := repository.NewClientRepository(dbPool)
	clientService := service.NewClientService(clientRepository)

	carRepository := repository.NewCarRepository(dbPool)
	carService := service.NewCarService(carRepository)

	employeeRepository := repository.NewEmployeeRepository(dbPool)

	employeePositionRepository :=
		repository.NewEmployeePositionRepository(dbPool)

	employeePhoneRepository :=
		repository.NewEmployeePhoneRepository(dbPool)

	employeeService := service.NewEmployeeService(
		employeeRepository,
		employeePositionRepository,
		employeePhoneRepository,
	)

	clientHandler := handler.NewClientHandler(
		clientService,
		carService,
	)

	carHandler := handler.NewCarHandler(
		carService,
		clientService,
	)

	employeeHandler := handler.NewEmployeeHandler(
		employeeService,
	)

	mux := router.SetupRoutes(
		clientHandler,
		carHandler,
		employeeHandler,
		loginHandler,
	)

	addr := ":" + cfg.Server.Port

	fmt.Println("Server started on", addr)

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
