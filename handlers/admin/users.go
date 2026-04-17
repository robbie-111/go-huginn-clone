package admin

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	adminComponents "go-huginn-clone/components/admin/users"
	"go-huginn-clone/components/layouts"
	"go-huginn-clone/middleware"
	"go-huginn-clone/models"
)

func buildProps(w http.ResponseWriter, r *http.Request, title string) layouts.PageProps {
	return layouts.PageProps{
		Title:       title,
		CurrentUser: middleware.GetCurrentUser(r),
		Flash:       middleware.GetFlashes(w, r),
		CurrentPath: r.URL.Path,
	}
}

func buildPagination(total int, basePath string) models.Pagination {
	return models.Pagination{CurrentPage: 1, TotalPages: 1, TotalCount: total, PerPage: 30, BasePath: basePath}
}

func UsersIndex(w http.ResponseWriter, r *http.Request) {
	props := buildProps(w, r, "Users")
	users := models.MockUsers()
	pagination := buildPagination(len(users), "/admin/users")
	currentUser := middleware.GetCurrentUser(r)
	adminComponents.Index(props, users, pagination, currentUser).Render(r.Context(), w)
}

func UsersNew(w http.ResponseWriter, r *http.Request) {
	props := buildProps(w, r, "Create User")
	user := models.User{}
	adminComponents.NewPage(props, user, nil).Render(r.Context(), w)
}

func UsersEdit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	users := models.MockUsers()
	var user models.User
	for _, u := range users {
		if u.ID == id {
			user = u
			break
		}
	}
	currentUser := middleware.GetCurrentUser(r)
	props := buildProps(w, r, "Edit "+user.Username)
	adminComponents.EditPage(props, user, nil, currentUser).Render(r.Context(), w)
}

func UsersCreate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "User was successfully created.")
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

func UsersUpdate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "User was successfully updated.")
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

func UsersDestroy(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "User was successfully deleted.")
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

func UsersDeactivate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "User deactivated.")
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

func UsersActivate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "User activated.")
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

func UsersSwitchToUser(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Switched to user.")
	http.Redirect(w, r, "/", http.StatusFound)
}

func UsersSwitchBack(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Switched back to admin.")
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}
