package router

import (
	"net/http"

	"github.com/Varfa/GarageHub/internal/handler"
	"github.com/Varfa/GarageHub/internal/middleware"
)

func SetupRoutes(
	clientHandler *handler.ClientHandler,
	carHandler *handler.CarHandler,
	employeeHandler *handler.EmployeeHandler,
	warehouseHandler *handler.WarehouseHandler,
	loginHandler *handler.LoginHandler,
	orderHandler *handler.OrderHandler,
	setupHandler *handler.SetupHandler,
	authMiddleware *middleware.AuthMiddleware,
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	private := func(
		handlerFunc http.HandlerFunc,
	) http.Handler {
		return authMiddleware.RequireAuth(
			handlerFunc,
		)
	}

	// ---------------------------------------------------------
	// Static files
	// ---------------------------------------------------------

	static := http.FileServer(
		http.Dir("../frontend/static"),
	)

	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			static,
		),
	)

	// ---------------------------------------------------------
	// Public pages
	// ---------------------------------------------------------

	mux.HandleFunc(
		"/language",
		loginHandler.ChangeLanguage,
	)
	mux.Handle(
		"/logout",
		private(
			loginHandler.Logout,
		),
	)
	mux.HandleFunc(
		"/",
		handler.RootRedirectHandler,
	)

	mux.HandleFunc(
		"/setup",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if r.Method == http.MethodGet {
				setupHandler.SetupPage(
					w,
					r,
				)
				return
			}

			setupHandler.CreateOwner(
				w,
				r,
			)
		},
	)

	mux.HandleFunc(
		"/login",
		loginHandler.Login,
	)

	// ---------------------------------------------------------
	// Dashboard
	// ---------------------------------------------------------

	mux.Handle(
		"/dashboard",
		private(
			handler.DashboardHandler,
		),
	)

	// ---------------------------------------------------------
	// Clients
	// ---------------------------------------------------------

	mux.Handle(
		"/clients",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"clients.view",
				http.HandlerFunc(
					clientHandler.Clients,
				),
			),
		),
	)

	mux.Handle(
		"/clients/create",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"clients.create",
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
						if r.Method == http.MethodGet {
							clientHandler.CreatePage(
								w,
								r,
							)
							return
						}

						clientHandler.Create(
							w,
							r,
						)
					},
				),
			),
		),
	)

	mux.Handle(
		"/clients/delete",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"clients.delete",
				http.HandlerFunc(
					clientHandler.Delete,
				),
			),
		),
	)

	mux.Handle(
		"/clients/view",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"clients.view",
				http.HandlerFunc(
					clientHandler.View,
				),
			),
		),
	)

	mux.Handle(
		"/clients/update",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"clients.edit",
				http.HandlerFunc(
					clientHandler.Update,
				),
			),
		),
	)
	// ---------------------------------------------------------
	// Cars
	// ---------------------------------------------------------

	mux.Handle(
		"/cars",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"cars.view",
				http.HandlerFunc(
					carHandler.Cars,
				),
			),
		),
	)

	mux.Handle(
		"/cars/create",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"cars.create",
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
						if r.Method == http.MethodGet {
							carHandler.CreatePage(
								w,
								r,
							)
							return
						}

						carHandler.Create(
							w,
							r,
						)
					},
				),
			),
		),
	)

	mux.Handle(
		"/cars/delete",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"cars.delete",
				http.HandlerFunc(
					carHandler.Delete,
				),
			),
		),
	)

	mux.Handle(
		"/cars/view",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"cars.view",
				http.HandlerFunc(
					carHandler.View,
				),
			),
		),
	)

	mux.Handle(
		"/cars/update",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"cars.edit",
				http.HandlerFunc(
					carHandler.Update,
				),
			),
		),
	)

	mux.Handle(
		"/cars/change-owner",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"cars.edit",
				http.HandlerFunc(
					carHandler.ChangeOwner,
				),
			),
		),
	)
	// ---------------------------------------------------------
	// Employees
	// ---------------------------------------------------------

	mux.Handle(
		"/employees",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"employees.view",
				http.HandlerFunc(
					employeeHandler.Employees,
				),
			),
		),
	)

	mux.Handle(
		"/employees/archive",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"employees.view",
				http.HandlerFunc(
					employeeHandler.ArchivePage,
				),
			),
		),
	)

	mux.Handle(
		"/employees/view",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"employees.view",
				http.HandlerFunc(
					employeeHandler.View,
				),
			),
		),
	)

	mux.Handle(
		"/employees/update",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"employees.edit",
				http.HandlerFunc(
					employeeHandler.Update,
				),
			),
		),
	)

	mux.Handle(
		"/employees/archive-action",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"employees.archive",
				http.HandlerFunc(
					employeeHandler.Archive,
				),
			),
		),
	)

	mux.Handle(
		"/employees/restore",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"employees.archive",
				http.HandlerFunc(
					employeeHandler.Restore,
				),
			),
		),
	)

	mux.Handle(
		"/employees/create",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"employees.create",
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
						if r.Method == http.MethodGet {
							employeeHandler.CreatePage(
								w,
								r,
							)
							return
						}

						employeeHandler.Create(
							w,
							r,
						)
					},
				),
			),
		),
	)

	mux.Handle(
		"/employees/phones/create",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"employees.edit",
				http.HandlerFunc(
					employeeHandler.AddPhone,
				),
			),
		),
	)

	mux.Handle(
		"/employees/phones/primary",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"employees.edit",
				http.HandlerFunc(
					employeeHandler.SetPrimaryPhone,
				),
			),
		),
	)

	mux.Handle(
		"/employees/phones/delete",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"employees.edit",
				http.HandlerFunc(
					employeeHandler.DeletePhone,
				),
			),
		),
	)
	// ---------------------------------------------------------
	// Users
	// ---------------------------------------------------------

	mux.Handle(
		"/users/create",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"users.create",
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
						if r.Method == http.MethodGet {
							userHandler.CreatePage(
								w,
								r,
							)
							return
						}

						userHandler.Create(
							w,
							r,
						)
					},
				),
			),
		),
	)
	mux.Handle(
		"/settings/users",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"users.view",
				http.HandlerFunc(
					userHandler.Users,
				),
			),
		),
	)
	mux.Handle(
		"/settings/users/edit",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"users.edit",
				http.HandlerFunc(
					userHandler.EditPage,
				),
			),
		),
	)
	mux.Handle(
		"/settings/users/role",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"users.edit",
				http.HandlerFunc(
					userHandler.UpdateRole,
				),
			),
		),
	)
	mux.Handle(
		"/settings/users/active",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"users.deactivate",
				http.HandlerFunc(
					userHandler.SetActive,
				),
			),
		),
	)
	// ---------------------------------------------------------
	// Orders
	// ---------------------------------------------------------

	mux.Handle(
		"/orders",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"orders.view",
				http.HandlerFunc(
					orderHandler.Orders,
				),
			),
		),
	)

	mux.Handle(
		"/orders/view",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"orders.view",
				http.HandlerFunc(
					orderHandler.View,
				),
			),
		),
	)

	mux.Handle(
		"/orders/create",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"orders.create",
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
						if r.Method == http.MethodGet {
							orderHandler.CreatePage(
								w,
								r,
							)
							return
						}

						orderHandler.Create(
							w,
							r,
						)
					},
				),
			),
		),
	)

	mux.Handle(
		"/orders/status",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"orders.change_status",
				http.HandlerFunc(
					orderHandler.UpdateStatus,
				),
			),
		),
	)

	mux.Handle(
		"/orders/note",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"orders.add_note",
				http.HandlerFunc(
					orderHandler.AddNote,
				),
			),
		),
	)

	mux.Handle(
		"/orders/closed",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"orders.view",
				http.HandlerFunc(
					orderHandler.Closed,
				),
			),
		),
	)

	mux.Handle(
		"/orders/restore",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"orders.restore",
				http.HandlerFunc(
					orderHandler.Restore,
				),
			),
		),
	)

	mux.Handle(
		"/orders/employees/assign",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"orders.assign_employee",
				http.HandlerFunc(
					orderHandler.AssignEmployee,
				),
			),
		),
	)

	mux.Handle(
		"/orders/employees/unassign",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"orders.assign_employee",
				http.HandlerFunc(
					orderHandler.UnassignEmployee,
				),
			),
		),
	)

	mux.Handle(
		"/orders/delete",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"orders.delete",
				http.HandlerFunc(
					orderHandler.Delete,
				),
			),
		),
	)
	// ---------------------------------------------------------
	// Parts
	// ---------------------------------------------------------

	mux.Handle(
		"/parts",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"parts.view",
				http.HandlerFunc(
					handler.PartsHandler,
				),
			),
		),
	)
	// ---------------------------------------------------------
	// Warehouse
	// ---------------------------------------------------------

	mux.Handle(
		"/warehouse",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"warehouse.view",
				http.HandlerFunc(
					warehouseHandler.Warehouse,
				),
			),
		),
	)

	mux.Handle(
		"/warehouse/create",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"warehouse.manage",
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
						if r.Method == http.MethodGet {
							warehouseHandler.CreatePage(
								w,
								r,
							)
							return
						}

						warehouseHandler.Create(
							w,
							r,
						)
					},
				),
			),
		),
	)

	mux.Handle(
		"/warehouse/view",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"warehouse.view",
				http.HandlerFunc(
					warehouseHandler.View,
				),
			),
		),
	)

	mux.Handle(
		"/warehouse/update",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"warehouse.manage",
				http.HandlerFunc(
					warehouseHandler.Update,
				),
			),
		),
	)

	mux.Handle(
		"/warehouse/archive",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"warehouse.view",
				http.HandlerFunc(
					warehouseHandler.ArchivePage,
				),
			),
		),
	)

	mux.Handle(
		"/warehouse/archive-item",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"warehouse.writeoff",
				http.HandlerFunc(
					warehouseHandler.Archive,
				),
			),
		),
	)

	mux.Handle(
		"/warehouse/restore",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"warehouse.manage",
				http.HandlerFunc(
					warehouseHandler.Restore,
				),
			),
		),
	)
	// ---------------------------------------------------------
	// Reports
	// ---------------------------------------------------------

	mux.Handle(
		"/reports",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"reports.view",
				http.HandlerFunc(
					handler.ReportsHandler,
				),
			),
		),
	)
	// ---------------------------------------------------------
	// Settings
	// ---------------------------------------------------------

	mux.Handle(
		"/settings",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"settings.view",
				http.HandlerFunc(
					handler.SettingsHandler,
				),
			),
		),
	)
	mux.Handle(
		"/settings/roles",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"roles.view",
				http.HandlerFunc(
					roleHandler.Roles,
				),
			),
		),
	)
	mux.Handle(
		"/settings/roles/view",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"roles.view",
				http.HandlerFunc(
					roleHandler.View,
				),
			),
		),
	)
	mux.Handle(
		"/settings/roles/permissions",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"roles.edit",
				http.HandlerFunc(
					roleHandler.UpdatePermissions,
				),
			),
		),
	)
	mux.Handle(
		"/settings/roles/create",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"roles.create",
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
						if r.Method == http.MethodGet {
							roleHandler.CreatePage(
								w,
								r,
							)
							return
						}

						roleHandler.Create(
							w,
							r,
						)
					},
				),
			),
		),
	)
	mux.Handle(
		"/settings/roles/update",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"roles.edit",
				http.HandlerFunc(
					roleHandler.Update,
				),
			),
		),
	)
	mux.Handle(
		"/settings/roles/delete",
		authMiddleware.RequireAuth(
			authMiddleware.RequirePermission(
				"roles.delete",
				http.HandlerFunc(
					roleHandler.Delete,
				),
			),
		),
	)
	// ---------------------------------------------------------
	// System
	// ---------------------------------------------------------

	mux.HandleFunc(
		"/health",
		handler.HealthHandler,
	)

	// ---------------------------------------------------------
	// API
	// ---------------------------------------------------------

	mux.HandleFunc(
		"/api/v1/hello",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			_, _ = w.Write(
				[]byte("Hello from API 👋"),
			)
		},
	)

	return mux
}
