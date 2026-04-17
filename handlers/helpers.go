package handlers

import (
	"net/http"

	"go-huginn-clone/components/layouts"
	"go-huginn-clone/middleware"
	"go-huginn-clone/models"
)

// makeProps builds the common PageProps used by all handlers
func makeProps(w http.ResponseWriter, r *http.Request, title string) layouts.PageProps {
	return layouts.PageProps{
		Title:       title,
		CurrentUser: middleware.GetCurrentUser(r),
		Flash:       middleware.GetFlashes(w, r),
		CurrentPath: r.URL.Path,
	}
}

// redirect is a convenience wrapper
func redirect(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusFound)
}

// mockPagination builds a pagination struct for mock data
func mockPagination(total int, basePath string) models.Pagination {
	return models.Pagination{
		CurrentPage: 1,
		TotalPages:  1,
		TotalCount:  total,
		PerPage:     30,
		BasePath:    basePath,
	}
}
