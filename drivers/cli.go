package drivers

import (
	"fmt"
	"time"

	"github.com/ioneldev/telemetry"
)

type CLIDriver struct{}

func (d *CLIDriver) Write(entry telemetry.LogEntry) {
	// Start building the log message
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
	fmt.Print(logMessage)
}
