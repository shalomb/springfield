package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLog_CreatesFiles(t *testing.T) {
	tmpDir := t.TempDir()
	originalLogDir := LogDir
	LogDir = tmpDir
	defer func() { LogDir = originalLogDir }()

	testAgent := "test-agent"
	logMessage := "test message"
	err := Log(logMessage, "INFO", testAgent, "", "", nil, 0, nil)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	agentLogPath := filepath.Join(tmpDir, testAgent+".log")
	if _, err := os.Stat(agentLogPath); os.IsNotExist(err) {
		t.Errorf("expected %s to exist", agentLogPath)
	}

	centralLogPath := filepath.Join(tmpDir, "springfield.log")
	if _, err := os.Stat(centralLogPath); os.IsNotExist(err) {
		t.Errorf("expected %s to exist", centralLogPath)
	}

	data, err := os.ReadFile(agentLogPath)
	if err != nil {
		t.Fatalf("failed to read agent log: %v", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	if entry.Message != logMessage {
		t.Errorf("expected message %q, got %q", logMessage, entry.Message)
	}
}

func TestLog_WithExtraFields(t *testing.T) {
	tmpDir := t.TempDir()
	originalLogDir := LogDir
	LogDir = tmpDir
	defer func() { LogDir = originalLogDir }()

	extraData := map[string]interface{}{"key": "value"}
	tokenUsage := map[string]int{"total": 100}
	err := Log("msg", "DEBUG", "agent1", "EPIC-1", "Task 1", tokenUsage, 0.05, extraData)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	agentLogPath := filepath.Join(tmpDir, "agent1.log")
	data, err := os.ReadFile(agentLogPath)
	if err != nil {
		t.Fatalf("failed to read agent log: %v", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "DEBUG" {
		t.Errorf("expected level DEBUG, got %q", entry.Level)
	}

	if entry.Cost != 0.05 {
		t.Errorf("expected cost 0.05, got %f", entry.Cost)
	}

	if val, ok := entry.Data["key"]; !ok || val != "value" {
		t.Errorf("expected data.key to be 'value', got %v", entry.Data["key"])
	}
}

func TestLog_ErrorPaths(t *testing.T) {
	// Test failure to create directory
	originalLogDir := LogDir
	defer func() { LogDir = originalLogDir }()

	// Use a path that is actually a file to trigger MkdirAll error
	tmpFile, err := os.CreateTemp("", "notadir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	LogDir = tmpFile.Name()
	err = Log("msg", "INFO", "agent", "", "", nil, 0, nil)
	if err == nil {
		t.Error("expected error when LogDir is a file, got nil")
	}
}

func TestLog_AppendError(t *testing.T) {
	tmpDir := t.TempDir()
	originalLogDir := LogDir
	LogDir = tmpDir
	defer func() { LogDir = originalLogDir }()

	// Create a directory where the log file should be to trigger appendToFile error
	// For agent log
	err := os.Mkdir(filepath.Join(tmpDir, "agent.log"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = Log("msg", "INFO", "agent", "", "", nil, 0, nil)
	if err == nil {
		t.Error("expected error when agent log is a directory, got nil")
	}

	// Now for central log
	// First, remove the agent.log directory so it succeeds
	os.Remove(filepath.Join(tmpDir, "agent.log"))
	// Then create springfield.log as a directory
	err = os.Mkdir(filepath.Join(tmpDir, "springfield.log"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = Log("msg", "INFO", "agent", "", "", nil, 0, nil)
	if err == nil {
		t.Error("expected error when central log is a directory, got nil")
	}
}

func TestLog_Defaults(t *testing.T) {
	tmpDir := t.TempDir()
	originalLogDir := LogDir
	LogDir = tmpDir
	defer func() { LogDir = originalLogDir }()

	err := Log("msg", "", "", "", "", nil, 0, nil)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	agentLogPath := filepath.Join(tmpDir, "unknown.log")
	data, _ := os.ReadFile(agentLogPath)
	var entry Entry
	_ = json.Unmarshal(data, &entry)

	if entry.Level != "INFO" {
		t.Errorf("expected default level INFO, got %q", entry.Level)
	}
	if entry.Agent != "unknown" {
		t.Errorf("expected default agent unknown, got %q", entry.Agent)
	}
}

func TestLog_MarshalError(t *testing.T) {
	// Data containing something that can't be marshaled (e.g. a function)
	badData := map[string]interface{}{"fn": func() {}}
	err := Log("msg", "INFO", "agent", "", "", nil, 0, badData)
	if err == nil {
		t.Error("expected marshal error, got nil")
	}
}

func TestAppendToFile_WriteError(t *testing.T) {
	// /dev/full returns ENOSPC on write
	if _, err := os.Stat("/dev/full"); os.IsNotExist(err) {
		t.Skip("/dev/full not available")
	}

	err := appendToFile("/dev/full", []byte("data"))
	if err == nil {
		t.Error("expected error writing to /dev/full, got nil")
	}
}
