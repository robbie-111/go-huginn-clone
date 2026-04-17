package components

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/a-h/templ"
	"go-huginn-clone/models"
)

// TimeAgo returns a human-readable relative time string (Rails time_ago_in_words equivalent)
func TimeAgo(t time.Time) string {
	diff := time.Since(t)
	minutes := diff.Minutes()
	hours := diff.Hours()
	days := hours / 24

	switch {
	case minutes < 1:
		return "less than a minute"
	case minutes < 2:
		return "1 minute"
	case minutes < 45:
		return fmt.Sprintf("%d minutes", int(minutes))
	case minutes < 90:
		return "about 1 hour"
	case hours < 24:
		return fmt.Sprintf("about %d hours", int(hours))
	case hours < 42:
		return "1 day"
	case days < 30:
		return fmt.Sprintf("%d days", int(days))
	case days < 45:
		return "about 1 month"
	case days < 365:
		return fmt.Sprintf("%d months", int(days/30))
	case days < 548:
		return "about 1 year"
	default:
		return fmt.Sprintf("%d years", int(days/365))
	}
}

// Pluralize returns "1 item" or "N items"
func Pluralize(count int, singular string) string {
	if count == 1 {
		return fmt.Sprintf("1 %s", singular)
	}
	// simple pluralization: add 's' unless already ends in 's'
	plural := singular
	if !strings.HasSuffix(singular, "s") {
		plural = singular + "s"
	}
	return fmt.Sprintf("%d %s", count, plural)
}

// Truncate shortens a string to length, appending "..." if needed
func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// PrettyJSON returns a pretty-printed JSON string
func PrettyJSON(v interface{}) string {
	if v == nil {
		return "{}"
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}

// YesNo returns "Yes" or "No"
func YesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// IconTag returns an HTML icon element for Bootstrap 3 glyphicon or Font Awesome
func IconTag(icon string, extraClasses ...string) templ.Component {
	var cls string
	extra := strings.Join(extraClasses, " ")

	switch {
	case strings.HasPrefix(icon, "fa-"):
		// Font Awesome
		cls = "fa-solid " + icon
		if extra != "" {
			cls += " " + extra
		}
	case strings.HasPrefix(icon, "glyphicon-"):
		// Bootstrap 3 Glyphicon
		cls = "glyphicon " + icon
		if extra != "" {
			cls += " " + extra
		}
	default:
		cls = "glyphicon glyphicon-" + icon
		if extra != "" {
			cls += " " + extra
		}
	}

	return templ.Raw(fmt.Sprintf(`<span class="%s" aria-hidden="true"></span>`, cls))
}

// AgentScheduleName returns a human-readable schedule name
func AgentScheduleName(schedule string) string {
	schedules := map[string]string{
		"never":     "Never",
		"every_1m":  "Every 1 minute",
		"every_2m":  "Every 2 minutes",
		"every_5m":  "Every 5 minutes",
		"every_10m": "Every 10 minutes",
		"every_30m": "Every 30 minutes",
		"every_1h":  "Every 1 hour",
		"every_2h":  "Every 2 hours",
		"every_6h":  "Every 6 hours",
		"every_12h": "Every 12 hours",
		"every_1d":  "Every 1 day",
		"every_2d":  "Every 2 days",
		"every_7d":  "Every 7 days",
		"midnight":  "Midnight",
		"1am":       "1am", "2am": "2am", "3am": "3am", "4am": "4am",
		"5am": "5am", "6am": "6am", "7am": "7am", "8am": "8am",
		"9am": "9am", "10am": "10am", "11am": "11am", "noon": "Noon",
	}
	if name, ok := schedules[schedule]; ok {
		return name
	}
	return schedule
}

// WorkingStatus returns an HTML badge for agent working status
func WorkingStatus(agent models.Agent) templ.Component {
	if agent.Disabled {
		return templ.Raw(`<span class="label label-default">Disabled</span>`)
	}
	if agent.Working {
		return templ.Raw(`<span class="label label-success">Yes</span>`)
	}
	return templ.Raw(`<span class="label label-danger">No</span>`)
}

// AgentTypeIcon returns a Font Awesome icon for the agent type
func AgentTypeIcon(agent models.Agent) templ.Component {
	iconMap := map[string]string{
		"WeatherAgent":       "fa-cloud-sun",
		"RssFeedAgent":       "fa-rss",
		"EmailAgent":         "fa-envelope",
		"EmailDigestAgent":   "fa-envelope-open",
		"TwitterStreamAgent": "fa-twitter",
		"WebhookAgent":       "fa-code",
		"DataOutputAgent":    "fa-database",
		"SlackAgent":         "fa-slack",
		"HttpRequestAgent":   "fa-globe",
		"JavaScriptAgent":    "fa-js",
		"TriggerAgent":       "fa-bolt",
		"PostAgent":          "fa-paper-plane",
		"ManualEventAgent":   "fa-hand-pointer",
	}

	shortType := agent.ShortType
	// strip "Agents::" prefix if present
	if idx := strings.LastIndex(agent.Type, "::"); idx >= 0 {
		shortType = agent.Type[idx+2:]
	}

	icon, ok := iconMap[shortType]
	if !ok {
		icon = "fa-robot"
	}

	cls := "fa-solid " + icon
	if agent.Disabled {
		cls += " text-muted"
	}

	return templ.Raw(fmt.Sprintf(`<span class="%s" title="%s"></span>`, cls, agent.ShortType))
}

// ScenarioLabel returns a colored badge for a scenario
func ScenarioLabel(scenario models.Scenario) templ.Component {
	bgColor := scenario.TagBgColor
	if bgColor == "" {
		bgColor = "#5bc0de"
	}
	fgColor := scenario.TagFgColor
	if fgColor == "" {
		fgColor = "#ffffff"
	}
	icon := scenario.Icon
	if icon == "" {
		icon = "gear"
	}
	html := fmt.Sprintf(
		`<span class="label scenario" style="background-color:%s;color:%s;"><span class="fa-solid fa-%s"></span> %s</span>`,
		bgColor, fgColor, icon, scenario.Name,
	)
	return templ.Raw(html)
}

// ServiceLabel returns a colored badge for an OAuth service
func ServiceLabel(service *models.Service) templ.Component {
	if service == nil {
		return templ.Raw(`<span class="label label-default">None</span>`)
	}
	colorMap := map[string]string{
		"twitter":  "#55acee",
		"github":   "#444444",
		"tumblr":   "#2c4762",
		"dropbox":  "#007EE5",
		"evernote": "#00C85A",
		"google":   "#4285F4",
	}
	iconMap := map[string]string{
		"twitter":  "fa-twitter",
		"github":   "fa-github",
		"tumblr":   "fa-tumblr",
		"dropbox":  "fa-dropbox",
		"evernote": "fa-leaf",
		"google":   "fa-google",
	}
	color := colorMap[service.Provider]
	if color == "" {
		color = "#777"
	}
	icon := iconMap[service.Provider]
	if icon == "" {
		icon = "fa-link"
	}
	html := fmt.Sprintf(
		`<span class="label label-service service-%s" style="background-color:%s;color:#fff;"><span class="fa-brands %s"></span> %s</span>`,
		service.Provider, color, icon, service.Name,
	)
	return templ.Raw(html)
}

// SortableColumn renders a sortable column header link
func SortableColumn(currentSort, currentDir, col, defaultDir, name string, path string) templ.Component {
	icon := ""
	nextDir := defaultDir
	if currentSort == col {
		if currentDir == "asc" {
			icon = ` <span class="fa-solid fa-sort-up"></span>`
			nextDir = "desc"
		} else {
			icon = ` <span class="fa-solid fa-sort-down"></span>`
			nextDir = "asc"
		}
	}
	html := fmt.Sprintf(
		`<a href="%s?sort=%s&direction=%s">%s%s</a>`,
		path, col, nextDir, name, icon,
	)
	return templ.Raw(html)
}

// Paginate renders Bootstrap 3 pagination
func Paginate(p models.Pagination) templ.Component {
	if p.TotalPages <= 1 {
		return templ.Raw("")
	}

	var sb strings.Builder
	sb.WriteString(`<nav><ul class="pagination">`)

	// Previous
	if p.CurrentPage > 1 {
		sb.WriteString(fmt.Sprintf(
			`<li><a href="%s?page=%d">&laquo;</a></li>`,
			p.BasePath, p.CurrentPage-1,
		))
	} else {
		sb.WriteString(`<li class="disabled"><span>&laquo;</span></li>`)
	}

	// Page numbers (show up to 7 around current)
	start := int(math.Max(1, float64(p.CurrentPage-3)))
	end := int(math.Min(float64(p.TotalPages), float64(p.CurrentPage+3)))

	if start > 1 {
		sb.WriteString(fmt.Sprintf(`<li><a href="%s?page=1">1</a></li>`, p.BasePath))
		if start > 2 {
			sb.WriteString(`<li class="disabled"><span>...</span></li>`)
		}
	}

	for i := start; i <= end; i++ {
		if i == p.CurrentPage {
			sb.WriteString(fmt.Sprintf(`<li class="active"><span>%d</span></li>`, i))
		} else {
			sb.WriteString(fmt.Sprintf(`<li><a href="%s?page=%d">%d</a></li>`, p.BasePath, i, i))
		}
	}

	if end < p.TotalPages {
		if end < p.TotalPages-1 {
			sb.WriteString(`<li class="disabled"><span>...</span></li>`)
		}
		sb.WriteString(fmt.Sprintf(`<li><a href="%s?page=%d">%d</a></li>`, p.BasePath, p.TotalPages, p.TotalPages))
	}

	// Next
	if p.CurrentPage < p.TotalPages {
		sb.WriteString(fmt.Sprintf(
			`<li><a href="%s?page=%d">&raquo;</a></li>`,
			p.BasePath, p.CurrentPage+1,
		))
	} else {
		sb.WriteString(`<li class="disabled"><span>&raquo;</span></li>`)
	}

	sb.WriteString(`</ul></nav>`)
	return templ.Raw(sb.String())
}

// OmniauthProviderIcon returns the FA icon for an OAuth provider
func OmniauthProviderIcon(provider string) templ.Component {
	iconMap := map[string]string{
		"twitter":  "fa-brands fa-twitter",
		"github":   "fa-brands fa-github",
		"google":   "fa-brands fa-google",
		"tumblr":   "fa-brands fa-tumblr",
		"dropbox":  "fa-brands fa-dropbox",
		"evernote": "fa-solid fa-leaf",
		"slack":    "fa-brands fa-slack",
	}
	icon := iconMap[provider]
	if icon == "" {
		icon = "fa-solid fa-link"
	}
	return templ.Raw(fmt.Sprintf(`<span class="%s"></span>`, icon))
}

// OmniauthProviderName returns a display name for an OAuth provider
func OmniauthProviderName(provider string) string {
	names := map[string]string{
		"twitter":  "Twitter / X",
		"github":   "GitHub",
		"google":   "Google",
		"tumblr":   "Tumblr",
		"dropbox":  "Dropbox",
		"evernote": "Evernote",
		"slack":    "Slack",
	}
	if name, ok := names[provider]; ok {
		return name
	}
	return strings.Title(provider)
}

// JobStatusLabel returns a colored label for a job status
func JobStatusLabel(job models.Job) templ.Component {
	switch job.Status {
	case "pending":
		return templ.Raw(`<span class="label label-default">Pending</span>`)
	case "queued":
		return templ.Raw(`<span class="label label-primary">Queued</span>`)
	case "running":
		return templ.Raw(`<span class="label label-info">Running</span>`)
	case "locked":
		return templ.Raw(`<span class="label label-warning">Locked</span>`)
	case "failed":
		return templ.Raw(`<span class="label label-danger">Failed</span>`)
	case "succeeded":
		return templ.Raw(`<span class="label label-success">Succeeded</span>`)
	default:
		return templ.Raw(fmt.Sprintf(`<span class="label label-default">%s</span>`, job.Status))
	}
}

// LogLevelIcon returns an icon for a log level
func LogLevelIcon(level int) templ.Component {
	switch {
	case level >= 4:
		return templ.Raw(`<span class="fa-solid fa-circle-xmark text-danger" title="Error"></span>`)
	case level == 3:
		return templ.Raw(`<span class="fa-solid fa-circle-xmark text-danger" title="Error"></span>`)
	case level == 2:
		return templ.Raw(`<span class="fa-solid fa-triangle-exclamation text-warning" title="Warning"></span>`)
	default:
		return templ.Raw(`<span class="fa-solid fa-circle-info text-info" title="Info"></span>`)
	}
}

// UserAccountState returns a label for user account state
func UserAccountState(user models.User) templ.Component {
	if !user.Active {
		return templ.Raw(`<span class="label label-danger">Deactivated</span>`)
	}
	if user.Admin {
		return templ.Raw(`<span class="label label-warning">Admin</span>`)
	}
	return templ.Raw(`<span class="label label-success">Active</span>`)
}

// SafeTimeDiff returns a relative duration string for future times
func SafeTimeDiff(t time.Time) string {
	diff := time.Until(t)
	if diff < 0 {
		return "overdue"
	}
	minutes := diff.Minutes()
	hours := diff.Hours()
	switch {
	case minutes < 1:
		return "in less than a minute"
	case minutes < 60:
		return fmt.Sprintf("in %d minutes", int(minutes))
	case hours < 24:
		return fmt.Sprintf("in %d hours", int(hours))
	default:
		return fmt.Sprintf("in %d days", int(hours/24))
	}
}

// KeepEventsForLabel returns a human-readable label for event retention
func KeepEventsForLabel(seconds int) string {
	schedules := models.MockEventRetentionSchedules()
	s := fmt.Sprintf("%d", seconds)
	for _, pair := range schedules {
		if pair[1] == s {
			return pair[0]
		}
	}
	return fmt.Sprintf("%d seconds", seconds)
}

// AgentFromJob returns the agent name for a job (mock)
func AgentFromJob(job models.Job) string {
	if job.AgentID != nil {
		if job.AgentName != "" {
			return job.AgentName
		}
		return fmt.Sprintf("Agent #%d", *job.AgentID)
	}
	if job.JobClass != "" {
		return job.JobClass
	}
	return "(system)"
}

// HighlightedEvent checks if an event id matches the highlight parameter
func HighlightedEvent(eventID int, hlParam string) bool {
	return fmt.Sprintf("%d", eventID) == hlParam
}

// OmniauthProviders returns the list of configured OAuth providers (mock)
func OmniauthProviders() []string {
	return []string{"twitter", "github", "google", "dropbox", "tumblr"}
}

// ServiceAgentCount returns the number of agents using a service (mock)
func ServiceAgentCount(serviceID int) int {
	counts := map[int]int{1: 2, 2: 0, 3: 1}
	return counts[serviceID]
}
