package telemetry

import (
	"os"
	"strings"
	"testing"
)

// Sample YAML content for testing
const sampleYAML = `
defaultLogLevel: 1
`

// TestLoadConfig tests the LoadConfig function
func TestLoadConfig(t *testing.T) {
	// Create a temporary file for the test
	tmpFile, err := os.CreateTemp("", "config_test.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temp file after the test

	// Write sample YAML content to the temp file
	if _, err := tmpFile.Write([]byte(sampleYAML)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	expectedLogLevel := LogLevel(1) // Adjust this based on your LogLevel type
	if config.DefaultLogLevel != expectedLogLevel {
		t.Errorf("expected DefaultLogLevel = %v, got %v", expectedLogLevel, config.DefaultLogLevel)
	}
}

// TestLoadConfigFileNotFound tests the LoadConfig function when the file is not found
func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := LoadConfig("non_existent_file.yaml")

	if err == nil {
		t.Fatalf("expected an error when loading a non-existent file, got nil")
	}

	expectedErrorMsg := "error reading config file: open non_existent_file.yaml: no such file or directory"
	if err.Error() != expectedErrorMsg {
		t.Errorf("expected file not found error, got %v", err)
	}
}

// TestLoadConfigUnmarshalError tests the LoadConfig function when unmarshaling fails
func TestLoadConfigUnmarshalError(t *testing.T) {
	// Create a temporary file with invalid YAML content
	invalidYAML := `
	defaultLogLevel: "not_a_number"  // This should be an integer
	`
	tmpFile, err := os.CreateTemp("", "invalid_config_test.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temp file after the test

	// Write invalid YAML content to the temp file
	if _, err := tmpFile.Write([]byte(invalidYAML)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close() // Close the file to ensure it's written

	// Load the config from the temporary file
	_, err = LoadConfig(tmpFile.Name())
	if err == nil {
		t.Fatalf("expected an error when unmarshaling invalid YAML, got nil")
	}

	// Optionally, you can check for a specific error message or type
	expectedErrorMsg := "error unmarshaling config" // Adjust this based on your actual error handling
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("expected error message to contain %q, got %v", expectedErrorMsg, err)
	}
}
