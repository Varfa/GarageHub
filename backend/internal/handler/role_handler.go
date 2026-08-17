package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/service"
)

type RoleCreatePageData struct {
	Error string
}
type RoleViewPageData struct {
	Role             models.RoleDetails
	PermissionGroups []models.PermissionGroup
	Error            string
}
type RoleHandler struct {
	roleService *service.RoleService
}

type RolesPageData struct {
	Roles []models.RoleListItem
	Error string
}

func NewRoleHandler(
	roleService *service.RoleService,
) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}
func (h *RoleHandler) Roles(
	w http.ResponseWriter,
	r *http.Request,
) {
	roles, err := h.roleService.List(
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

	data := RolesPageData{
		Roles: roles,
		Error: "",
	}

	RenderTemplate(
		w,
		r,
		"roles",
		data,
	)
}
func (h *RoleHandler) View(
	w http.ResponseWriter,
	r *http.Request,
) {
	roleID, err := strconv.ParseInt(
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

	role, err := h.roleService.GetByID(
		r.Context(),
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

	permissions, err := h.roleService.ListPermissions(
		r.Context(),
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
	var permissionGroups []models.PermissionGroup

	for _, permission := range permissions {
		if len(permissionGroups) == 0 ||
			permissionGroups[len(permissionGroups)-1].Module != permission.Module {

			permissionGroups = append(
				permissionGroups,
				models.PermissionGroup{
					Module: permission.Module,
				},
			)
		}

		lastGroup := &permissionGroups[len(permissionGroups)-1]

		lastGroup.Permissions = append(
			lastGroup.Permissions,
			permission,
		)
	}
	data := RoleViewPageData{
		Role:             *role,
		PermissionGroups: permissionGroups,
		Error:            r.URL.Query().Get("error")}

	RenderTemplate(
		w,
		r,
		"role_view",
		data,
	)
}
func (h *RoleHandler) UpdatePermissions(
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

	permissionValues := r.Form["permission_id"]

	permissionIDs := make(
		[]int64,
		0,
		len(permissionValues),
	)

	for _, value := range permissionValues {
		permissionID, err := strconv.ParseInt(
			value,
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

		permissionIDs = append(
			permissionIDs,
			permissionID,
		)
	}

	err = h.roleService.UpdatePermissions(
		r.Context(),
		roleID,
		permissionIDs,
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
		"/settings/roles/view?id="+strconv.FormatInt(roleID, 10),
		http.StatusSeeOther,
	)
}
func (h *RoleHandler) CreatePage(
	w http.ResponseWriter,
	r *http.Request,
) {
	data := RoleCreatePageData{
		Error: "",
	}

	RenderTemplate(
		w,
		r,
		"role_create",
		data,
	)
}
func (h *RoleHandler) Create(
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

	code := r.FormValue("code")
	name := r.FormValue("name")
	description := r.FormValue("description")

	roleID, err := h.roleService.Create(
		r.Context(),
		code,
		name,
		description,
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
		"/settings/roles/view?id="+strconv.FormatInt(roleID, 10),
		http.StatusSeeOther,
	)
}
func (h *RoleHandler) Update(
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

	name := r.FormValue("name")
	description := r.FormValue("description")
	isActive := r.FormValue("is_active") == "true"

	err = h.roleService.Update(
		r.Context(),
		roleID,
		name,
		description,
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
		"/settings/roles/view?id="+strconv.FormatInt(roleID, 10),
		http.StatusSeeOther,
	)
}
func (h *RoleHandler) Delete(
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

	err = h.roleService.Delete(
		r.Context(),
		roleID,
	)
	if err != nil {
		if errors.Is(
			err,
			service.ErrRoleInUse,
		) {
			http.Redirect(
				w,
				r,
				"/settings/roles/view?id="+
					strconv.FormatInt(roleID, 10)+
					"&error=role_in_use",
				http.StatusSeeOther,
			)
			return
		}

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
		"/settings/roles",
		http.StatusSeeOther,
	)
}
