package telemetry

import (
	"testing"
)

// MockDriver is a mock implementation of the Driver interface for testing
type MockDriver struct {
	entries []LogEntry
}

func (m *MockDriver) Write(entry LogEntry) {
	m.entries = append(m.entries, entry)
}

// ClearEntries clears the log entries in the mock driver
func (m *MockDriver) ClearEntries() {
	m.entries = []LogEntry{}
}

func TestDefaultLogger(t *testing.T) {
	mockDriver := &MockDriver{}
	logger := NewDefaultLogger(Debug)
	logger.setDrivers([]Driver{mockDriver})

	// Test logging without a transaction
	logger.Info("Info message", nil)
	if len(mockDriver.entries) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(mockDriver.entries))
	}
	if mockDriver.entries[0].Message != "Info message" {
		t.Errorf("Expected message 'Info message', got '%s'", mockDriver.entries[0].Message)
	}
	if mockDriver.entries[0].TransactionID != "" {
		t.Errorf("Expected empty TransactionID, got '%s'", mockDriver.entries[0].TransactionID)
	}

	mockDriver.ClearEntries()

	// Test starting a transaction
	logger.startTransaction("txn1", nil)
	logger.Info("Info message in transaction", nil)
	if len(mockDriver.entries) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(mockDriver.entries))
	}
	if mockDriver.entries[1].TransactionID != "txn1" {
		t.Errorf("Expected TransactionID 'txn1', got '%s'", mockDriver.entries[1].TransactionID)
	}

	mockDriver.ClearEntries()

	// Test ending a transaction
	logger.endTransaction("txn1")
	if len(mockDriver.entries) != 1 {
		t.Errorf("Expected 1 log entries after ending transaction, got %d", len(mockDriver.entries))
	}
	if mockDriver.entries[0].Message != "Ended transaction txn1" {
		t.Errorf("Expected message 'Ended transaction txn1', got '%s'", mockDriver.entries[0].Message)
	}
	logger.Info("Info message outside transaction", nil)
	if mockDriver.entries[1].TransactionID != "" {
		t.Errorf("Expected empty TransactionID after ending transaction, got '%s'", mockDriver.entries[1].TransactionID)
	}

	mockDriver.ClearEntries()

	// Test ending a non-existent transaction
	logger.endTransaction("non-existent")
	if len(mockDriver.entries) != 1 {
		t.Errorf("Expected 1 log entries after attempting to end non-existent transaction, got %d", len(mockDriver.entries))
	}
	if mockDriver.entries[0].Message != "Attempted to end non-existent transaction non-existent" {
		t.Errorf("Expected message 'Attempted to end non-existent transaction non-existent', got '%s'", mockDriver.entries[0].Message)
	}
}
