package drivers

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ioneldev/telemetry"
)

type TextDriver struct {
	filePath string
	mu       sync.Mutex
}

func NewTextDriver(filePath string) *TextDriver {
	return &TextDriver{filePath: filePath}
}

func (d *TextDriver) Write(entry telemetry.LogEntry) {
	d.mu.Lock()         // Lock to ensure thread safety
	defer d.mu.Unlock() // Unlock after writing

	// Open the text file
	file, err := os.OpenFile(d.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close() // Ensure the file is closed after writing

	// Write the log entry to the file
	logMessage := fmt.Sprintf("[%s] %s",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level)

	// Conditionally append TransactionID if it exists
	if entry.TransactionID != "" {
		logMessage += fmt.Sprintf(" (Transaction: %s)", entry.TransactionID)
	}

	logMessage += fmt.Sprintf(": %s", entry.Message)

	// Conditionally append Tags if they are defined
	if len(entry.Tags) > 0 {
		logMessage += fmt.Sprintf(" Tags: %v", entry.Tags)
	}

	// Print the final log message
	fmt.Println()
	fmt.Fprintln(file, logMessage)
}
