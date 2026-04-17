package models

import "time"

// User represents a Huginn user
type User struct {
	ID       int
	Username string
	Email    string
	Admin    bool
	Active   bool

	AgentCount         int
	InactiveAgentCount int
	EventCount         int
	RecentEventCount   int
	ScenarioCount      int
	CreatedAt          time.Time
}

// Agent represents a Huginn agent
type Agent struct {
	ID        int
	Name      string
	Type      string // e.g. "WeatherAgent", "TwitterAgent"
	ShortType string // e.g. "Weather Agent"
	Schedule  string
	Disabled  bool
	UserID    int

	CanBeScheduled        bool
	CanCreateEvents       bool
	CanReceiveEvents      bool
	CanControlOtherAgents bool
	CanDryRun             bool

	LastCheckAt   *time.Time
	LastEventAt   *time.Time
	LastReceiveAt *time.Time
	EventsCount   int
	KeepEventsFor int // seconds

	Working        bool
	WorkingMessage string
	RecentErrors   bool

	Options map[string]interface{}
	Memory  map[string]interface{}

	Sources        []Agent
	Receivers      []Agent
	Controllers    []Agent
	ControlTargets []Agent
	Scenarios      []Scenario

	SourceIDs        []int
	ReceiverIDs      []int
	ControllerIDs    []int
	ControlTargetIDs []int
	ScenarioIDs      []int

	PropagateImmediately bool
	Service              *Service

	Errors []string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a Agent) ShortTypeTitleized() string {
	return a.ShortType
}

func (a Agent) Unavailable() bool {
	return a.Disabled
}

// Event represents a Huginn event
type Event struct {
	ID        int
	AgentID   int
	Agent     *Agent
	Payload   map[string]interface{}
	Lat       *float64
	Lng       *float64
	ExpiresAt *time.Time
	CreatedAt time.Time
}

// Scenario represents a Huginn scenario
type Scenario struct {
	ID          int
	Name        string
	Description string
	Public      bool
	Icon        string
	TagBgColor  string
	TagFgColor  string
	SourceURL   string
	Agents      []Agent
	UserID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Job represents a background job
type Job struct {
	ID        int
	AgentID   *int
	AgentName string
	JobClass  string
	Status    string // "pending", "running", "failed", "succeeded"
	Attempts  int
	LastError string
	LockedAt  *time.Time
	LockedBy  string
	FailedAt  *time.Time
	RunAt     time.Time
	CreatedAt time.Time
}

// AgentLog represents a log entry for an agent
type AgentLog struct {
	ID              int
	AgentID         int
	Message         string
	Level           int // 0=debug,1=info,2=warn,3=error,4=fatal
	InboundEventID  *int
	OutboundEventID *int
	CreatedAt       time.Time
}

// UserCredential represents a stored credential
type UserCredential struct {
	ID              int
	UserID          int
	CredentialName  string
	CredentialValue string
	Mode            string // "text", "json"
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Service represents an OAuth service connection
type Service struct {
	ID        int
	UserID    int
	Provider  string // "twitter", "github", "google", etc.
	Name      string // OAuth username
	Global    bool
	CreatedAt time.Time
}

// AgentTypeInfo holds metadata about an agent type (for the new/edit form)
type AgentTypeInfo struct {
	Name                  string
	Description           string
	CanBeScheduled        bool
	CanCreateEvents       bool
	CanReceiveEvents      bool
	CanControlOtherAgents bool
	CanDryRun             bool
	Options               map[string]interface{}
}

// FlashMessage is used to pass flash notices between pages
type FlashMessage struct {
	Type    string // "notice", "alert"
	Message string
}

// Pagination holds pagination metadata
type Pagination struct {
	CurrentPage int
	TotalPages  int
	TotalCount  int
	PerPage     int
	BasePath    string
	Params      map[string]string
}

// DryRunResult holds the result of a dry run
type DryRunResult struct {
	Events []map[string]interface{}
	Log    string
	Memory map[string]interface{}
}

// ScenarioImport holds state for a scenario import process
type ScenarioImport struct {
	Step             int // 1 or 2
	URL              string
	Data             string
	ParsedData       map[string]interface{}
	ExistingScenario *Scenario
	Dangerous        bool
	Errors           []string
}

// WorkerStatus holds the current worker/job status
type WorkerStatus struct {
	Pending        int
	AwaitingRetry  int
	RecentFailures int
	EventsSince    int
	SinceID        *int
}
