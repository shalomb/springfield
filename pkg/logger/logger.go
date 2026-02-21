package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single log entry in JSON format.
type Entry struct {
	Timestamp  string                 `json:"timestamp"`
	Message    string                 `json:"message"`
	Level      string                 `json:"level"`
	Agent      string                 `json:"agent"`
	Epic       string                 `json:"epic,omitempty"`
	Task       string                 `json:"task,omitempty"`
	TokenUsage interface{}            `json:"token_usage,omitempty"`
	Cost       float64                `json:"cost,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

var LogDir = "logs"

func init() {
	if dir := os.Getenv("SPRINGFIELD_LOG_DIR"); dir != "" {
		LogDir = dir
	}
}

// Log writes a message to the agent's log file and the central springfield log file.
func Log(message string, level string, agent string, epic string, task string, tokenUsage interface{}, cost float64, data map[string]interface{}) error {
	entry := Entry{
		Timestamp:  time.Now().Format(time.RFC3339),
		Message:    message,
		Level:      level,
		Agent:      agent,
		Epic:       epic,
		Task:       task,
		TokenUsage: tokenUsage,
		Cost:       cost,
		Data:       data,
	}

	if entry.Level == "" {
		entry.Level = "INFO"
	}
	if entry.Agent == "" {
		entry.Agent = "unknown"
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// Ensure logs directory exists
	if err := os.MkdirAll(LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Write to agent-specific log
	agentLogPath := filepath.Join(LogDir, fmt.Sprintf("%s.log", entry.Agent))
	if err := appendToFile(agentLogPath, jsonData); err != nil {
		return err
	}

	// Write to central log
	centralLogPath := filepath.Join(LogDir, "springfield.log")
	if err := appendToFile(centralLogPath, jsonData); err != nil {
		return err
	}

	return nil
}

func appendToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", path, err)
	}
	defer func() {
		_ = f.Close()
	}()

	// Append newline and write in a single call for better atomicity in concurrent environments
	fullData := append(data, '\n')
	if _, err := f.Write(fullData); err != nil {
		return fmt.Errorf("failed to write to log file %s: %w", path, err)
	}

	return nil
}
