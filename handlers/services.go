package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	serviceComponents "go-huginn-clone/components/services"
	"go-huginn-clone/middleware"
	"go-huginn-clone/models"
)

func ServicesIndex(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Services")
	services := models.MockServices()
	counts := map[int]int{1: 2, 2: 0, 3: 1}
	serviceComponents.Index(props, services, counts, "provider", "asc").Render(r.Context(), w)
}

func ServicesDestroy(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Service disconnected.")
	redirect(w, r, "/services")
}

func ServicesToggleAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	_, _ = strconv.Atoi(idStr)
	middleware.SetFlash(w, r, "notice", "Service visibility updated.")
	redirect(w, r, "/services")
}
