package orchestrator

import (
	"testing"
)

func TestTDClient_RobustUnmarshal(t *testing.T) {
	tc := &TDClient{}

	// Case 1: Single object
	data1 := []byte(`{"id": "td-123", "type": "epic"}`)
	var issue1 Issue
	if err := tc.robustUnmarshal(data1, &issue1); err != nil {
		t.Fatalf("robustUnmarshal failed for single object: %v", err)
	}
	if issue1.ID != "td-123" {
		t.Errorf("expected ID td-123, got %s", issue1.ID)
	}

	// Case 2: Array of one object
	data2 := []byte(`[{"id": "td-456", "type": "epic"}]`)
	var issue2 Issue
	if err := tc.robustUnmarshal(data2, &issue2); err != nil {
		t.Fatalf("robustUnmarshal failed for array of one: %v", err)
	}
	if issue2.ID != "td-456" {
		t.Errorf("expected ID td-456, got %s", issue2.ID)
	}

	// Case 3: Empty array (should fail for single target)
	data3 := []byte(`[]`)
	var issue3 Issue
	if err := tc.robustUnmarshal(data3, &issue3); err == nil {
		t.Error("expected robustUnmarshal to fail for empty array when target is single object")
	}

	// Case 4: Array of multiple (should fail for single target)
	data4 := []byte(`[{"id": "1"}, {"id": "2"}]`)
	var issue4 Issue
	if err := tc.robustUnmarshal(data4, &issue4); err == nil {
		t.Error("expected robustUnmarshal to fail for multiple objects when target is single object")
	}
}
