package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	jobComponents "go-huginn-clone/components/jobs"
	"go-huginn-clone/middleware"
	"go-huginn-clone/models"
)

func JobsIndex(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Jobs")
	jobs := models.MockJobs()
	pagination := mockPagination(len(jobs), "/jobs")
	jobComponents.Index(props, jobs, pagination).Render(r.Context(), w)
}

func JobsDestroy(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Job deleted.")
	redirect(w, r, "/jobs")
}

func JobsRun(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "id")
	middleware.SetFlash(w, r, "notice", "Job queued for immediate execution.")
	redirect(w, r, "/jobs")
}

func JobsDestroyFailed(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Failed jobs removed.")
	redirect(w, r, "/jobs")
}

func JobsDestroyAll(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "All jobs removed.")
	redirect(w, r, "/jobs")
}

func JobsRetryQueued(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Queued jobs retried.")
	redirect(w, r, "/jobs")
}
