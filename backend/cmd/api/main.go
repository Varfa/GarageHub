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
	"github.com/Varfa/GarageHub/internal/middleware"
	"github.com/Varfa/GarageHub/internal/repository"
	"github.com/Varfa/GarageHub/internal/router"
	"github.com/Varfa/GarageHub/internal/service"
)

func main() {
	// ---------------------------------------------------------
	// Config
	// ---------------------------------------------------------

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// ---------------------------------------------------------
	// Database
	// ---------------------------------------------------------

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

	// ---------------------------------------------------------
	// i18n
	// ---------------------------------------------------------

	translator, err := i18n.NewManager(
		[]string{"en", "lt", "uk", "ru"},
		"en",
	)
	if err != nil {
		log.Fatal(err)
	}

	handler.SetTranslator(translator)

	// ---------------------------------------------------------
	// Users / Auth
	// ---------------------------------------------------------

	userRepository := repository.NewUserRepository(
		dbPool,
	)

	userService := service.NewUserService(
		userRepository,
	)

	sessionRepository := repository.NewSessionRepository(
		dbPool,
	)

	sessionService := service.NewSessionService(
		sessionRepository,
	)
	authMiddleware := middleware.NewAuthMiddleware(
		sessionService,
		userService,
	)

	setupHandler := handler.NewSetupHandler(
		userService,
	)

	loginHandler := handler.NewLoginHandler(
		translator,
		userService,
		sessionService,
	)
	// ---------------------------------------------------------
	// Roles
	// ---------------------------------------------------------

	roleRepository := repository.NewRoleRepository(
		dbPool,
	)

	roleService := service.NewRoleService(
		roleRepository,
	)
	roleHandler := handler.NewRoleHandler(
		roleService,
	)

	// ---------------------------------------------------------
	// Clients
	// ---------------------------------------------------------

	clientRepository := repository.NewClientRepository(
		dbPool,
	)

	clientService := service.NewClientService(
		clientRepository,
	)

	// ---------------------------------------------------------
	// Cars
	// ---------------------------------------------------------

	carRepository := repository.NewCarRepository(
		dbPool,
	)

	carService := service.NewCarService(
		carRepository,
	)

	// ---------------------------------------------------------
	// Employees
	// ---------------------------------------------------------

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
	userHandler := handler.NewUserHandler(
		userService,
		roleService,
		employeeService,
	)
	// ---------------------------------------------------------
	// Warehouse
	// ---------------------------------------------------------

	warehouseRepository := repository.NewWarehouseRepository(
		dbPool,
	)

	warehouseService := service.NewWarehouseService(
		warehouseRepository,
	)

	// ---------------------------------------------------------
	// Orders
	// ---------------------------------------------------------

	orderRepository := repository.NewOrderRepository(
		dbPool,
	)

	orderService := service.NewOrderService(
		orderRepository,
	)

	// Заметки механиков к заказу.
	orderNoteRepository := repository.NewOrderNoteRepository(
		dbPool,
	)

	orderNoteService := service.NewOrderNoteService(
		orderNoteRepository,
	)

	// Назначение сотрудников на заказы.
	orderEmployeeRepository := repository.NewOrderEmployeeRepository(
		dbPool,
	)

	orderEmployeeService := service.NewOrderEmployeeService(
		orderEmployeeRepository,
	)

	// ---------------------------------------------------------
	// Handlers
	// ---------------------------------------------------------

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

	orderHandler := handler.NewOrderHandler(
		orderService,
		clientService,
		carService,
		orderNoteService,
		orderEmployeeService,
		employeeService,
	)

	// ---------------------------------------------------------
	// Router
	// ---------------------------------------------------------

	mux := router.SetupRoutes(
		clientHandler,
		carHandler,
		employeeHandler,
		warehouseHandler,
		loginHandler,
		orderHandler,
		setupHandler,
		authMiddleware,
		userHandler,
		roleHandler,
	)

	// ---------------------------------------------------------
	// HTTP Server
	// ---------------------------------------------------------

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
