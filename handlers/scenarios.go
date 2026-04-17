package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	siComponents "go-huginn-clone/components/scenario_imports"
	scenarioComponents "go-huginn-clone/components/scenarios"
	"go-huginn-clone/middleware"
	"go-huginn-clone/models"
)

func ScenariosIndex(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Scenarios")
	scenarios := models.MockScenarios()
	pagination := mockPagination(len(scenarios), "/scenarios")
	scenarioComponents.Index(props, scenarios, pagination, "name", "asc").Render(r.Context(), w)
}

func ScenariosShow(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	scenario := models.MockScenario(id)
	props := makeProps(w, r, scenario.Name)
	pagination := mockPagination(len(scenario.Agents), fmt.Sprintf("/scenarios/%d", id))
	scenarioComponents.Show(props, *scenario, scenario.Agents, pagination, "name", "asc").Render(r.Context(), w)
}

func ScenariosNew(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Create Scenario")
	scenario := models.Scenario{TagBgColor: "#5bc0de", TagFgColor: "#ffffff"}
	scenarioComponents.NewPage(props, scenario).Render(r.Context(), w)
}

func ScenariosEdit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	scenario := models.MockScenario(id)
	props := makeProps(w, r, "Edit "+scenario.Name)
	scenarioComponents.EditPage(props, *scenario).Render(r.Context(), w)
}

func ScenariosCreate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Scenario was successfully created.")
	redirect(w, r, "/scenarios")
}

func ScenariosUpdate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Scenario was successfully updated.")
	idStr := chi.URLParam(r, "id")
	redirect(w, r, "/scenarios/"+idStr)
}

func ScenariosDestroy(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Scenario was successfully deleted.")
	redirect(w, r, "/scenarios")
}

func ScenariosShare(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	scenario := models.MockScenario(id)
	props := makeProps(w, r, "Share "+scenario.Name)
	exportURL := fmt.Sprintf("http://localhost:3000/scenarios/%d/export.json", id)
	scenarioComponents.SharePage(props, *scenario, exportURL).Render(r.Context(), w)
}

func ScenariosEnableDisableAgents(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Agents updated.")
	idStr := chi.URLParam(r, "id")
	redirect(w, r, "/scenarios/"+idStr)
}

func ScenarioImportsNew(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Import Scenario")
	si := models.ScenarioImport{Step: 1}
	siComponents.NewPage(props, si).Render(r.Context(), w)
}

func ScenarioImportsCreate(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Import Scenario - Step 2")
	si := models.ScenarioImport{Step: 2, Data: `{"name":"Imported Scenario","agents":[]}`, Dangerous: false}
	siComponents.NewPage(props, si).Render(r.Context(), w)
}
