package models

import (
	"fmt"
	"time"
)

// MockUser returns a sample logged-in user
func MockUser() *User {
	return &User{
		ID:                 1,
		Username:           "admin",
		Email:              "admin@example.com",
		Admin:              true,
		Active:             true,
		AgentCount:         5,
		InactiveAgentCount: 1,
		EventCount:         142,
		RecentEventCount:   23,
		ScenarioCount:      2,
		CreatedAt:          time.Now().Add(-30 * 24 * time.Hour),
	}
}

// MockAgents returns sample agents
func MockAgents() []Agent {
	now := time.Now()
	lastCheck := now.Add(-5 * time.Minute)
	lastEvent := now.Add(-2 * time.Hour)
	lastReceive := now.Add(-10 * time.Minute)

	scenario1 := Scenario{ID: 1, Name: "Morning Digest", TagBgColor: "#5bc0de", TagFgColor: "#ffffff", Icon: "cloud"}
	scenario2 := Scenario{ID: 2, Name: "Twitter Monitor", TagBgColor: "#d9534f", TagFgColor: "#ffffff", Icon: "twitter"}

	return []Agent{
		{
			ID:               1,
			Name:             "Weather Agent",
			Type:             "Agents::WeatherAgent",
			ShortType:        "Weather Agent",
			Schedule:         "every_1h",
			Disabled:         false,
			CanBeScheduled:   true,
			CanCreateEvents:  true,
			CanReceiveEvents: false,
			CanDryRun:        true,
			LastCheckAt:      &lastCheck,
			LastEventAt:      &lastEvent,
			EventsCount:      42,
			KeepEventsFor:    604800,
			Working:          true,
			WorkingMessage:   "Working",
			Scenarios:        []Scenario{scenario1},
			Options:          map[string]interface{}{"location": "Seoul", "api_key": "xxx"},
			Memory:           map[string]interface{}{},
			CreatedAt:        now.Add(-10 * 24 * time.Hour),
		},
		{
			ID:               2,
			Name:             "RSS Feed Agent",
			Type:             "Agents::RssFeedAgent",
			ShortType:        "Rss Feed Agent",
			Schedule:         "every_30m",
			Disabled:         false,
			CanBeScheduled:   true,
			CanCreateEvents:  true,
			CanReceiveEvents: false,
			CanDryRun:        true,
			LastCheckAt:      &lastCheck,
			LastEventAt:      &lastEvent,
			EventsCount:      128,
			KeepEventsFor:    604800,
			Working:          true,
			WorkingMessage:   "Working",
			Scenarios:        []Scenario{scenario1},
			Options:          map[string]interface{}{"url": "https://example.com/feed.xml"},
			Memory:           map[string]interface{}{},
			CreatedAt:        now.Add(-20 * 24 * time.Hour),
		},
		{
			ID:               3,
			Name:             "Email Digest Agent",
			Type:             "Agents::EmailDigestAgent",
			ShortType:        "Email Digest Agent",
			Schedule:         "every_1d",
			Disabled:         false,
			CanBeScheduled:   true,
			CanCreateEvents:  false,
			CanReceiveEvents: true,
			CanDryRun:        true,
			LastCheckAt:      &lastCheck,
			LastReceiveAt:    &lastReceive,
			EventsCount:      0,
			KeepEventsFor:    0,
			Working:          true,
			WorkingMessage:   "Working",
			Scenarios:        []Scenario{scenario1},
			SourceIDs:        []int{1, 2},
			Options:          map[string]interface{}{"subject": "Morning Digest", "to": "user@example.com"},
			Memory:           map[string]interface{}{},
			CreatedAt:        now.Add(-15 * 24 * time.Hour),
		},
		{
			ID:               4,
			Name:             "Twitter Stream Agent",
			Type:             "Agents::TwitterStreamAgent",
			ShortType:        "Twitter Stream Agent",
			Schedule:         "",
			Disabled:         false,
			CanBeScheduled:   false,
			CanCreateEvents:  true,
			CanReceiveEvents: false,
			CanDryRun:        false,
			LastEventAt:      &lastEvent,
			EventsCount:      512,
			KeepEventsFor:    604800,
			Working:          true,
			WorkingMessage:   "Working",
			Scenarios:        []Scenario{scenario2},
			Options:          map[string]interface{}{"keywords": []string{"golang", "huginn"}},
			Memory:           map[string]interface{}{},
			CreatedAt:        now.Add(-5 * 24 * time.Hour),
		},
		{
			ID:               5,
			Name:             "Data Output Agent",
			Type:             "Agents::DataOutputAgent",
			ShortType:        "Data Output Agent",
			Schedule:         "",
			Disabled:         true,
			CanBeScheduled:   false,
			CanCreateEvents:  false,
			CanReceiveEvents: true,
			CanDryRun:        true,
			Working:          false,
			WorkingMessage:   "Disabled",
			Scenarios:        []Scenario{},
			Options:          map[string]interface{}{"secrets": []string{"my-secret"}},
			Memory:           map[string]interface{}{},
			CreatedAt:        now.Add(-2 * 24 * time.Hour),
		},
	}
}

// MockAgent returns a single agent by ID
func MockAgent(id int) *Agent {
	agents := MockAgents()
	for i := range agents {
		if agents[i].ID == id {
			a := agents[i]
			// populate sources/receivers for show page
			if len(a.SourceIDs) > 0 {
				a.Sources = []Agent{agents[0], agents[1]}
			}
			a.Memory = map[string]interface{}{"last_run": "2024-01-01", "count": 42}
			return &a
		}
	}
	return &agents[0]
}

// MockEvents returns sample events
func MockEvents(agentID int) []Event {
	now := time.Now()
	agents := MockAgents()

	events := []Event{}
	for i := 1; i <= 10; i++ {
		agentIdx := (i - 1) % len(agents)
		if agentID > 0 {
			found := false
			for j := range agents {
				if agents[j].ID == agentID {
					agentIdx = j
					found = true
					break
				}
			}
			if !found {
				agentIdx = 0
			}
		}
		a := agents[agentIdx]
		events = append(events, Event{
			ID:        i,
			AgentID:   a.ID,
			Agent:     &a,
			Payload:   map[string]interface{}{"title": fmt.Sprintf("Event #%d", i), "url": "https://example.com", "score": i * 10},
			CreatedAt: now.Add(-time.Duration(i) * time.Hour),
		})
	}
	return events
}

// MockEvent returns a single event
func MockEvent(id int) *Event {
	events := MockEvents(0)
	for i := range events {
		if events[i].ID == id {
			return &events[i]
		}
	}
	return &events[0]
}

// MockScenarios returns sample scenarios
func MockScenarios() []Scenario {
	now := time.Now()
	agents := MockAgents()
	return []Scenario{
		{
			ID:          1,
			Name:        "Morning Digest",
			Description: "Collects weather and RSS feeds, then emails a daily digest every morning.",
			Public:      true,
			Icon:        "cloud",
			TagBgColor:  "#5bc0de",
			TagFgColor:  "#ffffff",
			Agents:      agents[:3],
			CreatedAt:   now.Add(-30 * 24 * time.Hour),
		},
		{
			ID:          2,
			Name:        "Twitter Monitor",
			Description: "Monitors Twitter for keywords and saves results.",
			Public:      false,
			Icon:        "twitter",
			TagBgColor:  "#55acee",
			TagFgColor:  "#ffffff",
			Agents:      agents[3:4],
			CreatedAt:   now.Add(-10 * 24 * time.Hour),
		},
	}
}

// MockScenario returns a single scenario
func MockScenario(id int) *Scenario {
	scenarios := MockScenarios()
	for i := range scenarios {
		if scenarios[i].ID == id {
			return &scenarios[i]
		}
	}
	return &scenarios[0]
}

// MockJobs returns sample background jobs
func MockJobs() []Job {
	now := time.Now()
	agentID1 := 1
	agentID2 := 2
	agentID3 := 3
	failedAt := now.Add(-1 * time.Hour)
	lockedAt := now.Add(-2 * time.Minute)
	lockedAt2 := now.Add(-1 * time.Minute)

	return []Job{
		{
			ID:        1,
			AgentID:   &agentID1,
			AgentName: "Weather Agent",
			Status:    "pending",
			Attempts:  0,
			RunAt:     now.Add(5 * time.Minute),
			CreatedAt: now.Add(-10 * time.Minute),
		},
		{
			ID:        2,
			AgentID:   &agentID2,
			AgentName: "RSS Feed Agent",
			Status:    "running",
			LockedAt:  &lockedAt,
			LockedBy:  "worker-1",
			Attempts:  1,
			RunAt:     now.Add(-2 * time.Minute),
			CreatedAt: now.Add(-15 * time.Minute),
		},
		{
			ID:        3,
			AgentID:   nil,
			JobClass:  "AgentCleanupJob",
			Status:    "failed",
			FailedAt:  &failedAt,
			Attempts:  3,
			LastError: "Error: connection timeout\n  huginn/agent.rb:42\n  huginn/scheduler.rb:100",
			RunAt:     now.Add(-1 * time.Hour),
			CreatedAt: now.Add(-2 * time.Hour),
		},
		{
			ID:        4,
			AgentID:   &agentID3,
			AgentName: "Email Digest Agent",
			Status:    "queued",
			Attempts:  0,
			RunAt:     now.Add(15 * time.Minute),
			CreatedAt: now.Add(-3 * time.Minute),
		},
		{
			ID:        5,
			AgentID:   nil,
			JobClass:  "ScheduledEventCleanupJob",
			Status:    "locked",
			LockedAt:  &lockedAt2,
			LockedBy:  "worker-2",
			Attempts:  2,
			RunAt:     now.Add(-1 * time.Minute),
			CreatedAt: now.Add(-30 * time.Minute),
		},
	}
}

// MockLogs returns sample agent logs
func MockLogs(agentID int) []AgentLog {
	now := time.Now()
	inboundID := 3
	outboundID := 4
	return []AgentLog{
		{
			ID:        1,
			AgentID:   agentID,
			Message:   "Successfully fetched weather data for Seoul",
			Level:     1,
			CreatedAt: now.Add(-5 * time.Minute),
		},
		{
			ID:              2,
			AgentID:         agentID,
			Message:         "Received event and emitted new event with temperature data",
			Level:           1,
			InboundEventID:  &inboundID,
			OutboundEventID: &outboundID,
			CreatedAt:       now.Add(-35 * time.Minute),
		},
		{
			ID:        3,
			AgentID:   agentID,
			Message:   "Warning: API rate limit approaching (80% used)",
			Level:     2,
			CreatedAt: now.Add(-2 * time.Hour),
		},
		{
			ID:        4,
			AgentID:   agentID,
			Message:   "Error: Failed to connect to weather API after 3 retries\n  connection refused: api.weather.example.com:443",
			Level:     4,
			CreatedAt: now.Add(-6 * time.Hour),
		},
	}
}

// MockUserCredentials returns sample credentials
func MockUserCredentials() []UserCredential {
	now := time.Now()
	return []UserCredential{
		{ID: 1, CredentialName: "twitter_consumer_key", CredentialValue: "AbCdEfGhIjKlMnOpQrStUvWxYz123456", Mode: "text", CreatedAt: now.Add(-30 * 24 * time.Hour)},
		{ID: 2, CredentialName: "twitter_consumer_secret", CredentialValue: "super_secret_value_here_xxxxxxxxxxxxxxxx", Mode: "text", CreatedAt: now.Add(-30 * 24 * time.Hour)},
		{ID: 3, CredentialName: "user_full_name", CredentialValue: "John Doe", Mode: "text", CreatedAt: now.Add(-20 * 24 * time.Hour)},
		{ID: 4, CredentialName: "openai_api_key", CredentialValue: "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", Mode: "text", CreatedAt: now.Add(-5 * 24 * time.Hour)},
		{ID: 5, CredentialName: "smtp_config", CredentialValue: `{"host":"smtp.gmail.com","port":587}`, Mode: "json", CreatedAt: now.Add(-1 * 24 * time.Hour)},
	}
}

// MockServices returns sample OAuth services
func MockServices() []Service {
	now := time.Now()
	return []Service{
		{ID: 1, Provider: "twitter", Name: "@huginn_user", Global: false, CreatedAt: now.Add(-30 * 24 * time.Hour)},
		{ID: 2, Provider: "github", Name: "huginn-user", Global: true, CreatedAt: now.Add(-20 * 24 * time.Hour)},
		{ID: 3, Provider: "google", Name: "huginn.user@gmail.com", Global: false, CreatedAt: now.Add(-10 * 24 * time.Hour)},
	}
}

// MockUsers returns sample users (admin view)
func MockUsers() []User {
	now := time.Now()
	return []User{
		{ID: 1, Username: "admin", Email: "admin@example.com", Admin: true, Active: true, AgentCount: 5, InactiveAgentCount: 1, CreatedAt: now.Add(-60 * 24 * time.Hour)},
		{ID: 2, Username: "alice", Email: "alice@example.com", Admin: false, Active: true, AgentCount: 12, InactiveAgentCount: 0, CreatedAt: now.Add(-30 * 24 * time.Hour)},
		{ID: 3, Username: "bob", Email: "bob@example.com", Admin: false, Active: false, AgentCount: 3, InactiveAgentCount: 3, CreatedAt: now.Add(-15 * 24 * time.Hour)},
	}
}

// MockAgentTypes returns the list of available agent types for the new/edit form
func MockAgentTypes() []string {
	return []string{
		"Agents::DataOutputAgent",
		"Agents::EmailAgent",
		"Agents::EmailDigestAgent",
		"Agents::EventFormattingAgent",
		"Agents::HttpRequestAgent",
		"Agents::ImapFolderAgent",
		"Agents::JavaScriptAgent",
		"Agents::LiquidOutputAgent",
		"Agents::ManualEventAgent",
		"Agents::PeakDetectorAgent",
		"Agents::PostAgent",
		"Agents::RssFeedAgent",
		"Agents::SendGridAgent",
		"Agents::SlackAgent",
		"Agents::TriggerAgent",
		"Agents::TwitterStreamAgent",
		"Agents::WebhookAgent",
		"Agents::WeatherAgent",
	}
}

// MockSchedules returns available schedule options
func MockSchedules() [][]string {
	return [][]string{
		{"Never", "never"},
		{"Every 1 minute", "every_1m"},
		{"Every 2 minutes", "every_2m"},
		{"Every 5 minutes", "every_5m"},
		{"Every 10 minutes", "every_10m"},
		{"Every 30 minutes", "every_30m"},
		{"Every 1 hour", "every_1h"},
		{"Every 2 hours", "every_2h"},
		{"Every 6 hours", "every_6h"},
		{"Every 12 hours", "every_12h"},
		{"Every 1 day", "every_1d"},
		{"Every 2 days", "every_2d"},
		{"Every 7 days", "every_7d"},
		{"Midnight", "midnight"},
		{"1am", "1am"},
		{"2am", "2am"},
		{"3am", "3am"},
		{"4am", "4am"},
		{"5am", "5am"},
		{"6am", "6am"},
		{"7am", "7am"},
		{"8am", "8am"},
		{"9am", "9am"},
		{"10am", "10am"},
		{"11am", "11am"},
		{"Noon", "noon"},
	}
}

// MockEventRetentionSchedules returns event retention options
func MockEventRetentionSchedules() [][]string {
	return [][]string{
		{"Forever", "0"},
		{"1 minute", "60"},
		{"1 hour", "3600"},
		{"6 hours", "21600"},
		{"1 day", "86400"},
		{"3 days", "259200"},
		{"5 days", "432000"},
		{"7 days", "604800"},
		{"2 weeks", "1209600"},
		{"4 weeks", "2419200"},
	}
}

// MockPagination returns a sample pagination struct
func MockPagination(total, page, perPage int, basePath string) Pagination {
	totalPages := total / perPage
	if total%perPage != 0 {
		totalPages++
	}
	return Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  total,
		PerPage:     perPage,
		BasePath:    basePath,
	}
}
