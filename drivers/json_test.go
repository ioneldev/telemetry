package drivers

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/ioneldev/telemetry"
)

func TestWrite(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	driver := &JSONDriver{filePath: tempFile.Name()}

	// Test case: Successful write
	entry := telemetry.LogEntry{
		Level:         1,
		Message:       "Test message",
		TransactionID: "12345",
		Tags:          map[string]string{"tag1": "11", "tag2": "22"},
	}
	driver.Write(entry)

	// Read back the content to verify
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	var logEntry telemetry.LogEntry
	if err := json.Unmarshal(content, &logEntry); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if !reflect.DeepEqual(entry, logEntry) {
		t.Errorf("Expected log entry %v, got %v", entry, logEntry)
	}
}
