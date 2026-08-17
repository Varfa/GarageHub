package handler

import (
	"net/http"
	"strconv"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/service"
)

type UserEditPageData struct {
	User  models.UserListItem
	Roles []models.Role
	Error string
}

type UsersPageData struct {
	Users []models.UserListItem
	Error string
}
type UserCreatePageData struct {
	Employee models.Employee
	Roles    []models.Role
	Error    string
}
type UserHandler struct {
	userService     *service.UserService
	roleService     *service.RoleService
	employeeService *service.EmployeeService
}

func NewUserHandler(
	userService *service.UserService,
	roleService *service.RoleService,
	employeeService *service.EmployeeService,
) *UserHandler {
	return &UserHandler{
		userService:     userService,
		roleService:     roleService,
		employeeService: employeeService,
	}
}
func (h *UserHandler) CreatePage(
	w http.ResponseWriter,
	r *http.Request,
) {
	roles, err := h.roleService.ListActive(r.Context())
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	employeeID, err := strconv.Atoi(
		r.URL.Query().Get("id"))
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}
	employee, err := h.employeeService.GetByID(r.Context(), int64(employeeID))

	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	data := UserCreatePageData{
		Employee: *employee,
		Roles:    roles,
		Error:    "",
	}
	RenderTemplate(w, r, "user_create", data)

}
func (h *UserHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	employeeIDString := r.FormValue("employee_id")
	roleIDString := r.FormValue("role_id")
	email := r.FormValue("email")
	password := r.FormValue("password")

	employeeID, err := strconv.ParseInt(
		employeeIDString,
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	roleID, err := strconv.ParseInt(
		roleIDString,
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	err = h.userService.CreateUser(
		r.Context(),
		employeeID,
		roleID,
		email,
		password,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/employees/view?id="+strconv.FormatInt(employeeID, 10),
		http.StatusSeeOther,
	)
}
func (h *UserHandler) Users(
	w http.ResponseWriter,
	r *http.Request,
) {
	users, err := h.userService.List(r.Context())
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	data := UsersPageData{
		Users: users,
		Error: "",
	}
	RenderTemplate(
		w,
		r,
		"users",
		data,
	)
}
func (h *UserHandler) EditPage(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, err := strconv.ParseInt(
		r.URL.Query().Get("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	user, err := h.userService.GetByID(
		r.Context(),
		userID,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	if user.RoleID != nil {
		user.RoleIDValue = *user.RoleID
	}

	roles, err := h.roleService.ListActive(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	data := UserEditPageData{
		User:  *user,
		Roles: roles,
		Error: "",
	}

	RenderTemplate(
		w,
		r,
		"user_edit",
		data,
	)
}
func (h *UserHandler) UpdateRole(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	userID, err := strconv.ParseInt(
		r.FormValue("user_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	roleID, err := strconv.ParseInt(
		r.FormValue("role_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	err = h.userService.UpdateRole(
		r.Context(),
		userID,
		roleID,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/settings/users/edit?id="+strconv.FormatInt(userID, 10),
		http.StatusSeeOther,
	)
}
func (h *UserHandler) SetActive(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	userID, err := strconv.ParseInt(
		r.FormValue("user_id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	isActive := r.FormValue("is_active") == "true"

	err = h.userService.SetActive(
		r.Context(),
		userID,
		isActive,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/settings/users/edit?id="+strconv.FormatInt(userID, 10),
		http.StatusSeeOther,
	)
}
