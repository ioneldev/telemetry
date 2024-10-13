package telemetry

import (
	"fmt"
	"time"
)

type Driver interface {
	Write(entry LogEntry)
}

type LogEntry struct {
	Timestamp     time.Time         `json:"timestamp"`
	Level         LogLevel          `json:"level"`
	Message       string            `json:"message"`
	Tags          map[string]string `json:"tags"`
	TransactionID string            `json:"transaction_id"`
}

type Logger interface {
	Info(message string, tags map[string]string)
	Debug(message string, tags map[string]string)
	Warning(message string, tags map[string]string)
	Error(message string, tags map[string]string)

	log(entry LogEntry)
	startTransaction(transactionID string, attributes map[string]string)
	endTransaction(transactionID string)
	setDrivers(drivers []Driver)
}

type DefaultLogger struct {
	defaultLogLevel    LogLevel
	drivers            []Driver
	activeTransactions map[string]map[string]string
	currentTransaction string
}

func NewDefaultLogger(defaultLogLevel LogLevel) Logger {
	return &DefaultLogger{
		defaultLogLevel:    defaultLogLevel,
		drivers:            []Driver{},
		activeTransactions: make(map[string]map[string]string),
	}
}

// Ensure DefaultLogger implements Logger
var _ Logger = (*DefaultLogger)(nil)

func (l *DefaultLogger) setDrivers(drivers []Driver) {
	l.drivers = drivers
}

func (l *DefaultLogger) startTransaction(transactionID string, attributes map[string]string) {
	if _, exists := l.activeTransactions[transactionID]; exists {
		l.Warning(fmt.Sprintf("Transaction %s already exists. Overwriting.", transactionID), nil)
	}

	l.activeTransactions[transactionID] = attributes
	l.currentTransaction = transactionID

	l.Info(fmt.Sprintf("Started transaction %s", transactionID), attributes)
}

func (l *DefaultLogger) endTransaction(transactionID string) {
	if attributes, exists := l.activeTransactions[transactionID]; exists {
		l.Info(fmt.Sprintf("Ended transaction %s", transactionID), attributes)
		delete(l.activeTransactions, transactionID)
		if l.currentTransaction == transactionID {
			l.currentTransaction = ""
		}
	} else {
		l.Warning(fmt.Sprintf("Attempted to end non-existent transaction %s", transactionID), nil)
	}
}

func (l *DefaultLogger) log(entry LogEntry) {
	if len(l.drivers) == 0 {
		fmt.Println("Error logging message. There are no drivers set")
		return
	}

	if l.currentTransaction != "" {
		entry.TransactionID = l.currentTransaction
	}

	for _, driver := range l.drivers {
		driver.Write(entry)
	}
}

func (l *DefaultLogger) Debug(message string, tags map[string]string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     Debug,
		Message:   message,
		Tags:      tags,
	}

	l.log(entry)
}

func (l *DefaultLogger) Info(message string, tags map[string]string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     Info,
		Message:   message,
		Tags:      tags,
	}

	l.log(entry)
}

func (l *DefaultLogger) Warning(message string, tags map[string]string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     Warning,
		Message:   message,
		Tags:      tags,
	}

	l.log(entry)
}

func (l *DefaultLogger) Error(message string, tags map[string]string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     Error,
		Message:   message,
		Tags:      tags,
	}

	l.log(entry)
}
