package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	agentComponents "go-huginn-clone/components/agents"
	"go-huginn-clone/middleware"
	"go-huginn-clone/models"
)

func AgentsIndex(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Agents")
	agents := models.MockAgents()
	pagination := mockPagination(len(agents), "/agents")
	agentComponents.Index(props, agents, pagination, "name", "asc", nil, false).Render(r.Context(), w)
}

func AgentsNew(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Create Agent")
	props.LoadAceEditor = true
	props.LoadJSONEditor = true
	agent := models.Agent{
		Options: map[string]interface{}{},
		Memory:  map[string]interface{}{},
	}
	allAgents := models.MockAgents()
	scenarios := models.MockScenarios()
	agentComponents.New(props, agent, allAgents, models.MockAgentTypes(), models.MockSchedules(), models.MockEventRetentionSchedules(), scenarios).Render(r.Context(), w)
}

func AgentsShow(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	agent := models.MockAgent(id)
	if agent == nil {
		http.NotFound(w, r)
		return
	}
	props := makeProps(w, r, agent.Name)
	logs := models.MockLogs(agent.ID)
	agentComponents.Show(props, *agent, logs).Render(r.Context(), w)
}

func AgentsEdit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	agent := models.MockAgent(id)
	if agent == nil {
		http.NotFound(w, r)
		return
	}
	props := makeProps(w, r, "Edit "+agent.Name)
	props.LoadAceEditor = true
	props.LoadJSONEditor = true
	allAgents := models.MockAgents()
	scenarios := models.MockScenarios()
	agentComponents.Edit(props, *agent, allAgents, models.MockAgentTypes(), models.MockSchedules(), models.MockEventRetentionSchedules(), scenarios).Render(r.Context(), w)
}

func AgentsCreate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Agent was successfully created.")
	redirect(w, r, "/agents")
}

func AgentsUpdate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Agent was successfully updated.")
	idStr := chi.URLParam(r, "id")
	redirect(w, r, "/agents/"+idStr)
}

func AgentsDestroy(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Agent was successfully deleted.")
	redirect(w, r, "/agents")
}

func AgentsRun(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Agent run queued.")
	returnTo := r.URL.Query().Get("return")
	if returnTo == "" {
		returnTo = "/agents"
	}
	redirect(w, r, returnTo)
}

func AgentsPropagate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Event propagation queued.")
	redirect(w, r, "/agents")
}

func AgentsToggleVisibility(w http.ResponseWriter, r *http.Request) {
	redirect(w, r, "/agents")
}

func AgentsLeaveScenario(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Agent removed from scenario.")
	redirect(w, r, "/agents")
}

func AgentsReemitEvents(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Events re-emitted.")
	redirect(w, r, "/agents")
}

func AgentsRemoveEvents(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Events deleted.")
	redirect(w, r, "/agents")
}

func AgentsDestroyMemory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func AgentsTypeDetails(w http.ResponseWriter, r *http.Request) {
	agentType := r.URL.Query().Get("type")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"type":                     agentType,
		"can_be_scheduled":         true,
		"can_create_events":        true,
		"can_receive_events":       false,
		"can_control_other_agents": false,
		"can_dry_run":              true,
		"options":                  map[string]interface{}{},
		"description":              "<p>This agent type would display its description here.</p>",
	})
}

func AgentsValidate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func AgentsComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]string{})
}

func AgentsEventDescriptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func AgentsDryRunIndex(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "agent_id")
	id, _ := strconv.Atoi(idStr)
	agent := models.MockAgent(id)
	if agent == nil {
		http.NotFound(w, r)
		return
	}
	props := makeProps(w, r, "Dry Run - "+agent.Name)
	withEventMode := r.URL.Query().Get("with_event_mode")
	if withEventMode == "" {
		withEventMode = "no"
	}
	recentEvents := models.MockEvents(agent.ID)[:2]
	agentComponents.DryRunIndex(props, *agent, recentEvents, withEventMode).Render(r.Context(), w)
}

func AgentsDryRunCreate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "agent_id")
	id, _ := strconv.Atoi(idStr)
	agent := models.MockAgent(id)
	if agent == nil {
		http.NotFound(w, r)
		return
	}
	props := makeProps(w, r, "Dry Run Result - "+agent.Name)
	result := models.DryRunResult{
		Events: []map[string]interface{}{
			{"title": "Sample Event", "url": "https://example.com", "score": 42},
		},
		Log:    "[2024-01-01 12:00:00] INFO: Agent ran successfully\n[2024-01-01 12:00:00] INFO: Created 1 event",
		Memory: map[string]interface{}{"last_run": "2024-01-01"},
	}
	agentComponents.DryRunCreate(props, *agent, result).Render(r.Context(), w)
}

func AgentsLogsIndex(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "agent_id")
	id, _ := strconv.Atoi(idStr)
	logs := models.MockLogs(id)
	agentComponents.LogsTable(logs).Render(r.Context(), w)
}

func AgentsLogsClear(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
