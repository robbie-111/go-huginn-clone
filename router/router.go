package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-huginn-clone/handlers"
	adminHandlers "go-huginn-clone/handlers/admin"
	mw "go-huginn-clone/middleware"
)

func New() http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(mw.SetCurrentUser)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Public routes (no auth required)
	r.Get("/", handlers.HomeIndex)
	r.Get("/users/sign_in", handlers.AuthLoginPage)
	r.Post("/users/sign_in", handlers.AuthLogin)
	r.Delete("/users/sign_out", handlers.AuthLogout)
	r.Post("/users/sign_out", handlers.AuthLogout) // fallback for non-JS
	r.Get("/users/sign_up", handlers.AuthRegisterPage)
	r.Post("/users", handlers.AuthRegister)

	// Worker status (public JSON API)
	r.Get("/worker_status", handlers.WorkerStatus)

	// Auth-required routes
	r.Group(func(r chi.Router) {
		r.Use(mw.RequireLogin)

		// Account
		r.Get("/users/edit", handlers.AuthAccountEdit)
		r.Put("/users", handlers.AuthAccountUpdate)
		r.Patch("/users", handlers.AuthAccountUpdate)

		// Agents
		r.Get("/agents", handlers.AgentsIndex)
		r.Get("/agents/new", handlers.AgentsNew)
		r.Post("/agents", handlers.AgentsCreate)
		r.Put("/agents/toggle_visibility", handlers.AgentsToggleVisibility)
		r.Post("/agents/propagate", handlers.AgentsPropagate)
		r.Get("/agents/type_details", handlers.AgentsTypeDetails)
		r.Get("/agents/event_descriptions", handlers.AgentsEventDescriptions)
		r.Post("/agents/validate", handlers.AgentsValidate)
		r.Post("/agents/complete", handlers.AgentsComplete)
		r.Get("/agents/{id}", handlers.AgentsShow)
		r.Get("/agents/{id}/edit", handlers.AgentsEdit)
		r.Put("/agents/{id}", handlers.AgentsUpdate)
		r.Patch("/agents/{id}", handlers.AgentsUpdate)
		r.Delete("/agents/{id}", handlers.AgentsDestroy)
		r.Post("/agents/{id}/run", handlers.AgentsRun)
		r.Put("/agents/{id}/leave_scenario", handlers.AgentsLeaveScenario)
		r.Post("/agents/{id}/reemit_events", handlers.AgentsReemitEvents)
		r.Delete("/agents/{id}/remove_events", handlers.AgentsRemoveEvents)
		r.Delete("/agents/{id}/memory", handlers.AgentsDestroyMemory)

		// Agent logs
		r.Get("/agents/{agent_id}/logs", handlers.AgentsLogsIndex)
		r.Delete("/agents/{agent_id}/logs/clear", handlers.AgentsLogsClear)

		// Agent events
		r.Get("/agents/{agent_id}/events", handlers.AgentEventsIndex)

		// Dry runs
		r.Get("/agents/{agent_id}/dry_runs", handlers.AgentsDryRunIndex)
		r.Post("/agents/{agent_id}/dry_runs", handlers.AgentsDryRunCreate)
		r.Get("/dry_runs", handlers.AgentsDryRunIndex)
		r.Post("/dry_runs", handlers.AgentsDryRunCreate)

		// Events
		r.Get("/events", handlers.EventsIndex)
		r.Get("/events/{id}", handlers.EventsShow)
		r.Delete("/events/{id}", handlers.EventsDestroy)
		r.Post("/events/{id}/reemit", handlers.EventsReemit)

		// Diagram
		r.Get("/diagram", handlers.DiagramsShow)

		// Scenarios
		r.Get("/scenarios", handlers.ScenariosIndex)
		r.Get("/scenarios/new", handlers.ScenariosNew)
		r.Post("/scenarios", handlers.ScenariosCreate)
		r.Get("/scenarios/{id}", handlers.ScenariosShow)
		r.Get("/scenarios/{id}/edit", handlers.ScenariosEdit)
		r.Put("/scenarios/{id}", handlers.ScenariosUpdate)
		r.Patch("/scenarios/{id}", handlers.ScenariosUpdate)
		r.Delete("/scenarios/{id}", handlers.ScenariosDestroy)
		r.Get("/scenarios/{id}/share", handlers.ScenariosShare)
		r.Put("/scenarios/{id}/enable_or_disable_all_agents", handlers.ScenariosEnableDisableAgents)
		r.Get("/scenarios/{id}/diagram", handlers.ScenarioDiagramsShow)

		// Scenario imports
		r.Get("/scenario_imports/new", handlers.ScenarioImportsNew)
		r.Post("/scenario_imports", handlers.ScenarioImportsCreate)

		// User Credentials
		r.Get("/user_credentials", handlers.CredentialsIndex)
		r.Get("/user_credentials/new", handlers.CredentialsNew)
		r.Post("/user_credentials", handlers.CredentialsCreate)
		r.Post("/user_credentials/import", handlers.CredentialsImport)
		r.Get("/user_credentials/{id}/edit", handlers.CredentialsEdit)
		r.Put("/user_credentials/{id}", handlers.CredentialsUpdate)
		r.Patch("/user_credentials/{id}", handlers.CredentialsUpdate)
		r.Delete("/user_credentials/{id}", handlers.CredentialsDestroy)

		// Services
		r.Get("/services", handlers.ServicesIndex)
		r.Delete("/services/{id}", handlers.ServicesDestroy)
		r.Post("/services/{id}/toggle_availability", handlers.ServicesToggleAvailability)

		// Admin only
		r.Group(func(r chi.Router) {
			r.Use(mw.RequireAdmin)

			r.Get("/jobs", handlers.JobsIndex)
			r.Delete("/jobs/{id}", handlers.JobsDestroy)
			r.Put("/jobs/{id}/run", handlers.JobsRun)
			r.Delete("/jobs/destroy_failed", handlers.JobsDestroyFailed)
			r.Delete("/jobs/destroy_all", handlers.JobsDestroyAll)
			r.Post("/jobs/retry_queued", handlers.JobsRetryQueued)

			r.Get("/admin/users", adminHandlers.UsersIndex)
			r.Get("/admin/users/new", adminHandlers.UsersNew)
			r.Post("/admin/users", adminHandlers.UsersCreate)
			r.Get("/admin/users/switch_back", adminHandlers.UsersSwitchBack)
			r.Get("/admin/users/{id}/edit", adminHandlers.UsersEdit)
			r.Put("/admin/users/{id}", adminHandlers.UsersUpdate)
			r.Patch("/admin/users/{id}", adminHandlers.UsersUpdate)
			r.Delete("/admin/users/{id}", adminHandlers.UsersDestroy)
			r.Put("/admin/users/{id}/deactivate", adminHandlers.UsersDeactivate)
			r.Put("/admin/users/{id}/activate", adminHandlers.UsersActivate)
			r.Get("/admin/users/{id}/switch_to_user", adminHandlers.UsersSwitchToUser)
		})
	})

	return r
}
