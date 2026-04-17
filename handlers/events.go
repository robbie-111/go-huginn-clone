package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	eventComponents "go-huginn-clone/components/events"
	"go-huginn-clone/middleware"
	"go-huginn-clone/models"
)

func EventsIndex(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Events")
	events := models.MockEvents(0)
	pagination := mockPagination(len(events), "/events")
	hlParam := r.URL.Query().Get("hl")
	eventComponents.Index(props, events, pagination, nil, hlParam).Render(r.Context(), w)
}

func AgentEventsIndex(w http.ResponseWriter, r *http.Request) {
	agentIDStr := chi.URLParam(r, "agent_id")
	agentID, _ := strconv.Atoi(agentIDStr)
	agent := models.MockAgent(agentID)
	events := models.MockEvents(agentID)
	props := makeProps(w, r, agent.Name+"'s Events")
	pagination := mockPagination(len(events), r.URL.Path)
	hlParam := r.URL.Query().Get("hl")
	eventComponents.Index(props, events, pagination, agent, hlParam).Render(r.Context(), w)
}

func EventsShow(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	event := models.MockEvent(id)
	props := makeProps(w, r, "Event "+idStr)
	eventComponents.Show(props, *event).Render(r.Context(), w)
}

func EventsDestroy(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Event deleted.")
	redirect(w, r, "/events")
}

func EventsReemit(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Event re-emitted.")
	redirect(w, r, "/events")
}
