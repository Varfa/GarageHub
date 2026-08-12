package main

import (
	"context"
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
	// Запускаем миграции до старта приложения.
	// Если миграция не прошла — сервер не запускаем.
	if err := database.RunMigrations(
		context.Background(),
		dbPool,
	); err != nil {
		log.Fatal(err)
	}

	translator, err := i18n.NewManager(
		[]string{"en", "lt", "uk", "ru"},
		"en",
	)
	if err != nil {
		log.Fatal(err)
	}

	handler.SetTranslator(translator)

	loginHandler := handler.NewLoginHandler(
		translator,
	)

	// Clients
	clientRepository := repository.NewClientRepository(
		dbPool,
	)

	clientService := service.NewClientService(
		clientRepository,
	)

	// Cars
	carRepository := repository.NewCarRepository(
		dbPool,
	)

	carService := service.NewCarService(
		carRepository,
	)

	// Employees
	employeeRepository := repository.NewEmployeeRepository(
		dbPool,
	)

	employeePositionRepository :=
		repository.NewEmployeePositionRepository(
			dbPool,
		)

	employeePhoneRepository :=
		repository.NewEmployeePhoneRepository(
			dbPool,
		)

	employeeService := service.NewEmployeeService(
		employeeRepository,
		employeePositionRepository,
		employeePhoneRepository,
	)

	// Warehouse
	warehouseRepository := repository.NewWarehouseRepository(
		dbPool,
	)

	warehouseService := service.NewWarehouseService(
		warehouseRepository,
	)

	// Handlers
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

	warehouseHandler := handler.NewWarehouseHandler(
		warehouseService,
	)

	// Orders

	orderRepository := repository.NewOrderRepository(dbPool)
	orderService := service.NewOrderService(orderRepository)

	// Заметки механиков к заказу.
	orderNoteRepository := repository.NewOrderNoteRepository(dbPool)
	orderNoteService := service.NewOrderNoteService(orderNoteRepository)

	orderHandler := handler.NewOrderHandler(
		orderService,
		clientService,
		carService,
		orderNoteService,
	)

	// Router
	mux := router.SetupRoutes(
		clientHandler,
		carHandler,
		employeeHandler,
		warehouseHandler,
		loginHandler,
		orderHandler,
	)

	addr := ":" + cfg.Server.Port

	fmt.Println(
		"Server started on",
		addr,
	)

	err = http.ListenAndServe(
		addr,
		mux,
	)
	if err != nil {
		log.Fatal(err)
	}
}
