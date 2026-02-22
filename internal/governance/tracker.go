package governance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// UsageStats tracks token usage and cost.
type UsageStats struct {
	TotalTokens   int   `json:"total_tokens"`
	TotalCostNano int64 `json:"total_cost_nano"`
}

// DailyUsageMap maps dates (YYYY-MM-DD) to usage stats.
type DailyUsageMap map[string]UsageStats

// UsageTracker manages persistence of daily usage data.
type UsageTracker struct {
	LogDir string
	mu     sync.Mutex
}

// NewUsageTracker creates a new tracker using the specified log directory.
func NewUsageTracker(logDir string) *UsageTracker {
	return &UsageTracker{LogDir: logDir}
}

func (t *UsageTracker) getFilePath() string {
	return filepath.Join(t.LogDir, "usage.json")
}

// RecordUsage adds usage to the current day's stats.
func (t *UsageTracker) RecordUsage(tokens int, costNano int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	usage, err := t.loadUsage()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if usage == nil {
		usage = make(DailyUsageMap)
	}

	today := time.Now().Format("2006-01-02")
	stats := usage[today]
	stats.TotalTokens += tokens
	stats.TotalCostNano += costNano
	usage[today] = stats

	return t.saveUsage(usage)
}

// GetDailyUsage returns the usage stats for the current day.
func (t *UsageTracker) GetDailyUsage() (UsageStats, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	usage, err := t.loadUsage()
	if err != nil {
		if os.IsNotExist(err) {
			return UsageStats{}, nil
		}
		return UsageStats{}, err
	}

	today := time.Now().Format("2006-01-02")
	return usage[today], nil
}

func (t *UsageTracker) loadUsage() (DailyUsageMap, error) {
	data, err := os.ReadFile(t.getFilePath())
	if err != nil {
		return nil, err
	}

	var usage DailyUsageMap
	if err := json.Unmarshal(data, &usage); err != nil {
		return nil, err
	}

	return usage, nil
}

func (t *UsageTracker) saveUsage(usage DailyUsageMap) error {
	if err := os.MkdirAll(t.LogDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(usage, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(t.getFilePath(), data, 0644)
}
