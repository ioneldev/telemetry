package drivers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/ioneldev/telemetry"
)

// TestCLIDriver_Write tests the Write method of the CLIDriver
func TestCLIDriver_Write(t *testing.T) {
	driver := &CLIDriver{}

	entry := telemetry.LogEntry{
		Timestamp:     time.Now(),
		Level:         1,
		Message:       "Test message",
		TransactionID: "12345",
		Tags:          map[string]string{"tag1": "11", "tag2": "22"},
	}

	// Redirect os.Stdout so that we can check the output
	originalStdout := os.Stdout
	var buf bytes.Buffer
	r, w, _ := os.Pipe()
	os.Stdout = w
	driver.Write(entry)

	w.Close()
	output, _ := io.ReadAll(r)
	os.Stdout = originalStdout

	buf.Write(output) // Capture the output in the buffer

	// Check the output
	expectedOutput := fmt.Sprintf("[%s] %s (Transaction: %s): %s Tags: %v",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level,
		entry.TransactionID,
		entry.Message,
		entry.Tags)

	if string(output) != expectedOutput {
		t.Errorf("expected output %q, got %q", expectedOutput, output)
	}
}
