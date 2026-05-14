package models

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
)

const LogFilePath = "static/data/agent_logs.json"
const MaxLogsPerAgent = 50

var logMu sync.Mutex

// loadAllLogs reads all logs from the JSON file (must be called with logMu held)
func loadAllLogs() []AgentLog {
	data, err := os.ReadFile(LogFilePath)
	if err != nil || len(data) == 0 {
		return []AgentLog{}
	}
	var logs []AgentLog
	if err := json.Unmarshal(data, &logs); err != nil {
		return []AgentLog{}
	}
	return logs
}

// saveAllLogs writes all logs to the JSON file (must be called with logMu held)
func saveAllLogs(logs []AgentLog) error {
	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(LogFilePath, data, 0644)
}

// nextLogID returns the next available log ID (must be called with logMu held)
func nextLogID(logs []AgentLog) int {
	maxID := 0
	for _, l := range logs {
		if l.ID > maxID {
			maxID = l.ID
		}
	}
	return maxID + 1
}

// LoadLogs returns logs for the given agentID, sorted newest first
func LoadLogs(agentID int) []AgentLog {
	logMu.Lock()
	defer logMu.Unlock()

	all := loadAllLogs()
	var result []AgentLog
	for _, l := range all {
		if l.AgentID == agentID {
			result = append(result, l)
		}
	}
	// Sort newest first
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})
	return result
}

// AppendLog adds a log entry to the JSON file.
// If the agent already has MaxLogsPerAgent logs, the oldest one is removed.
func AppendLog(log AgentLog) error {
	logMu.Lock()
	defer logMu.Unlock()

	all := loadAllLogs()
	log.ID = nextLogID(all)

	// Count existing logs for this agent and remove oldest if over limit
	var agentLogs []AgentLog
	var otherLogs []AgentLog
	for _, l := range all {
		if l.AgentID == log.AgentID {
			agentLogs = append(agentLogs, l)
		} else {
			otherLogs = append(otherLogs, l)
		}
	}

	// Sort agent logs oldest first, trim to MaxLogsPerAgent-1
	sort.Slice(agentLogs, func(i, j int) bool {
		return agentLogs[i].CreatedAt.Before(agentLogs[j].CreatedAt)
	})
	if len(agentLogs) >= MaxLogsPerAgent {
		agentLogs = agentLogs[len(agentLogs)-MaxLogsPerAgent+1:]
	}

	agentLogs = append(agentLogs, log)
	all = append(otherLogs, agentLogs...)

	return saveAllLogs(all)
}

// ClearLogs removes all logs for the given agentID from the JSON file.
func ClearLogs(agentID int) error {
	logMu.Lock()
	defer logMu.Unlock()

	all := loadAllLogs()
	var remaining []AgentLog
	for _, l := range all {
		if l.AgentID != agentID {
			remaining = append(remaining, l)
		}
	}
	if remaining == nil {
		remaining = []AgentLog{}
	}
	return saveAllLogs(remaining)
}
