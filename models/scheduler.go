package models

import (
	"math/rand"
	"time"
)

// dummyLogEntries is the pool of dummy log messages to randomly pick from.
// Each entry: {level, message}
var dummyLogEntries = []struct {
	Level   int
	Message string
}{
	// Level 1 - Info (short)
	{1, "Successfully fetched weather data for Seoul"},
	{1, "RSS feed fetched: 12 new items found in https://feeds.example.com/tech"},
	{1, "Received event and emitted new event with temperature data"},
	{1, "HTTP GET https://api.example.com/v2/data completed in 342ms, status 200"},
	{1, "Email digest sent to user@example.com with 8 articles"},
	{1, "Agent check completed. No new events."},
	{1, "Webhook received and processed successfully. Payload size: 1.2KB"},
	{1, "Twitter stream connected. Tracking: 3 keywords"},

	// Level 2 - Warning (medium)
	{2, "Warning: API rate limit approaching (80% used). Remaining: 120 calls."},
	{2, "Slow response from upstream: 4821ms (threshold: 3000ms). Retrying with backoff."},
	{2, "Retrying request (attempt 2/3): connection timeout to api.example.com:443"},
	{2, "RSS feed returned 0 items. Feed may be empty or temporarily unavailable: https://feeds.example.com/news"},
	{2, "Memory usage high: 847MB / 1024MB. Consider clearing old events."},

	// Level 4 - Error (long, with stack traces / JSON payloads)
	{4, "Error: Failed to connect to weather API after 3 retries\n  connection refused: api.weather.example.com:443\n\nStack trace:\n  goroutine 23 [running]:\n  go-huginn-clone/agents.(*WeatherAgent).Check(0xc0001a4000)\n      /app/agents/weather_agent.go:87 +0x2c4\n  go-huginn-clone/scheduler.(*Runner).runAgent(0xc0002b8000, 0xc0001a4000)\n      /app/scheduler/runner.go:134 +0x98\n  created by go-huginn-clone/scheduler.(*Runner).Start\n      /app/scheduler/runner.go:56 +0x17c"},
	{4, "Unexpected JSON structure in API response:\n  Expected field 'temperature' not found in response body.\n\nReceived payload:\n{\n  \"status\": \"error\",\n  \"code\": 503,\n  \"message\": \"Service temporarily unavailable. Our engineers have been notified.\",\n  \"retry_after\": 60,\n  \"request_id\": \"req_8f3k2j9x1p\",\n  \"timestamp\": \"2026-05-13T09:14:22Z\"\n}\n\nFull URL: https://api.weather.example.com/v2/current?location=Seoul&units=metric&lang=en"},
	{4, "OAuth token expired for Twitter service.\n  Error: 401 Unauthorized\n\nResponse body:\n{\n  \"errors\": [\n    {\n      \"code\": 89,\n      \"message\": \"Invalid or expired token. Please re-authenticate.\",\n      \"label\": \"invalid_token\"\n    }\n  ]\n}\n\nAction required: Please reauthorize the service at /services\nToken issued at: 2026-03-01T00:00:00Z\nToken expired at: 2026-05-01T00:00:00Z"},
	{4, "Database connection pool exhausted after 30.2s wait:\n  Pool config: max_connections=10, min_idle=2, max_idle=5\n  Current state: active=10, idle=0, pending_requests=47\n\nLast 3 failed queries:\n  [1] SELECT * FROM events WHERE agent_id=3 AND created_at > '2026-05-12' ORDER BY created_at DESC LIMIT 100  (waited 30.1s)\n  [2] UPDATE agents SET last_check_at=NOW() WHERE id=3  (waited 30.2s)\n  [3] INSERT INTO agent_logs (agent_id, message, level) VALUES (3, ...)  (waited 30.2s)\n\nError: context deadline exceeded\nRecommendation: Increase pool size or reduce query frequency."},
	{4, "RssFeedAgent failed to parse feed content:\n  URL: https://feeds.example.com/tech-news\n  HTTP Status: 200 OK\n  Content-Type: text/html; charset=utf-8  (expected application/rss+xml or application/atom+xml)\n\nParser error: XML syntax error on line 1: invalid character '\\x00' in XML\n\nFirst 500 bytes of response:\n<!DOCTYPE html><html><head><title>503 Service Unavailable</title></head><body><h1>Service Unavailable</h1><p>The server is temporarily unable to service your request due to maintenance downtime or capacity problems. Please try again later.</p></body></html>\n\nThis usually means the feed server is returning an error page instead of the actual feed."},
	{4, "EmailDigestAgent SMTP delivery failed:\n  SMTP server: smtp.gmail.com:587\n  From: huginn@example.com\n  To: user@example.com\n  Subject: Morning Digest - May 13, 2026\n\nSMTP error: 550 5.1.1 The email account that you tried to reach does not exist.\n  Please try double-checking the recipient's email address for typos or unnecessary spaces.\n\nDelivery attempts: 3\nLast attempt: 2026-05-13T09:00:05Z\nNext retry: disabled (permanent failure)\n\nFull SMTP transcript:\n  -> 220 smtp.gmail.com ESMTP\n  <- EHLO huginn.example.com\n  -> 250-smtp.gmail.com at your service\n  <- AUTH LOGIN\n  -> 235 2.7.0 Accepted\n  <- MAIL FROM:<huginn@example.com>\n  -> 250 2.1.0 OK\n  <- RCPT TO:<user@example.com>\n  -> 550 5.1.1 user does not exist"},
}

// SeedInitialLogs seeds the JSON file with initial logs if it is empty.
// It uses the current MockAgents IDs (1–5) and the existing mock log patterns.
func SeedInitialLogs() {
	logMu.Lock()
	existing := loadAllLogs()
	logMu.Unlock()

	if len(existing) > 0 {
		return
	}

	now := time.Now()
	agentIDs := []int{1, 2, 3, 4, 5}

	seedEntries := []struct {
		agentID int
		level   int
		message string
		offset  time.Duration
	}{
		{1, 1, "Successfully fetched weather data for Seoul", -5 * time.Minute},
		{1, 1, "Received event and emitted new event with temperature data", -35 * time.Minute},
		{1, 2, "Warning: API rate limit approaching (80% used). Remaining: 120 calls.", -2 * time.Hour},
		{1, 4, "Error: Failed to connect to weather API after 3 retries\n  connection refused: api.weather.example.com:443\n\nStack trace:\n  goroutine 23 [running]:\n  go-huginn-clone/agents.(*WeatherAgent).Check(0xc0001a4000)\n      /app/agents/weather_agent.go:87 +0x2c4", -6 * time.Hour},
	}

	for _, e := range seedEntries {
		inboundID := 3
		outboundID := 4
		log := AgentLog{
			AgentID:   e.agentID,
			Message:   e.message,
			Level:     e.level,
			CreatedAt: now.Add(e.offset),
		}
		if e.agentID == 1 && e.level == 1 && e.offset == -35*time.Minute {
			log.InboundEventID = &inboundID
			log.OutboundEventID = &outboundID
		}
		_ = AppendLog(log)
	}

	// Seed a few logs for agents 2–5
	for _, agentID := range agentIDs[1:] {
		for i, entry := range dummyLogEntries[:4] {
			_ = AppendLog(AgentLog{
				AgentID:   agentID,
				Message:   entry.Message,
				Level:     entry.Level,
				CreatedAt: now.Add(-time.Duration(i+1) * 20 * time.Minute),
			})
		}
	}
}

// StartLogScheduler runs a background goroutine that appends a random dummy log
// for each agent (IDs 1–5) every interval.
func StartLogScheduler(interval time.Duration) {
	agentIDs := []int{1, 2, 3, 4, 5}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for _, agentID := range agentIDs {
			entry := dummyLogEntries[rng.Intn(len(dummyLogEntries))]
			log := AgentLog{
				AgentID:   agentID,
				Message:   entry.Message,
				Level:     entry.Level,
				CreatedAt: time.Now(),
			}
			_ = AppendLog(log)
		}
	}
}
