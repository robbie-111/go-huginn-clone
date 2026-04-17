package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	diagramComponents "go-huginn-clone/components/diagrams"
	"go-huginn-clone/models"
)

func DiagramsShow(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Agent Event Flow")
	props.LoadDiagram = true
	agents := models.MockAgents()
	diagramComponents.Show(props, agents, nil).Render(r.Context(), w)
}

func ScenarioDiagramsShow(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	scenario := models.MockScenario(id)
	props := makeProps(w, r, scenario.Name+"'s Agent Event Flow")
	props.LoadDiagram = true
	diagramComponents.Show(props, scenario.Agents, scenario).Render(r.Context(), w)
}
