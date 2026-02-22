package governance

import (
	"os"
	"testing"
)

func TestUsageTracker(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "usage-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tracker := NewUsageTracker(tmpDir)

	// Initial usage should be zero
	stats, err := tracker.GetDailyUsage()
	if err != nil {
		t.Fatalf("GetDailyUsage failed: %v", err)
	}
	if stats.TotalTokens != 0 || stats.TotalCostNano != 0 {
		t.Errorf("expected zero stats, got %+v", stats)
	}

	// Record some usage
	err = tracker.RecordUsage(100, 5000)
	if err != nil {
		t.Fatalf("RecordUsage failed: %v", err)
	}

	// Check updated usage
	stats, err = tracker.GetDailyUsage()
	if err != nil {
		t.Fatalf("GetDailyUsage failed: %v", err)
	}
	if stats.TotalTokens != 100 || stats.TotalCostNano != 5000 {
		t.Errorf("expected 100/5000, got %+v", stats)
	}

	// Record more usage
	err = tracker.RecordUsage(200, 10000)
	if err != nil {
		t.Fatalf("RecordUsage failed: %v", err)
	}

	// Check aggregated usage
	stats, err = tracker.GetDailyUsage()
	if err != nil {
		t.Fatalf("GetDailyUsage failed: %v", err)
	}
	if stats.TotalTokens != 300 || stats.TotalCostNano != 15000 {
		t.Errorf("expected 300/15000, got %+v", stats)
	}
}
