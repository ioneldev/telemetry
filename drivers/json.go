package drivers

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ioneldev/telemetry"
)

type JSONDriver struct {
	filePath string
	mu       sync.Mutex
}

func NewJSONDriver(filePath string) *JSONDriver {
	return &JSONDriver{filePath: filePath}
}

func (d *JSONDriver) Write(entry telemetry.LogEntry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	file, err := os.OpenFile(d.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("Error marshaling entry to JSON: %v\n", err)
		return
	}

	if _, err := file.Write(jsonEntry); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	if _, err := file.WriteString("\n"); err != nil {
		fmt.Printf("Error writing newline to file: %v\n", err)
	}
}
