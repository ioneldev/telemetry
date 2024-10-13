package telemetry

import (
	"fmt"
	"sync"
)

// Telemetry is a struct that holds the configuration, logger, drivers, and a mutex for thread safety.
// It provides methods for logging messages, starting and ending transactions, and managing drivers.
type Telemetry struct {
	config  Config
	logger  Logger
	drivers []Driver
	mu      sync.RWMutex
}

// New creates a new Telemetry instance with the provided config path.
// The config path should point to a valid configuration file.
// Returns an error if the configuration file cannot be loaded.
func New(configPath string) (*Telemetry, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	logger := NewDefaultLogger(config.DefaultLogLevel)

	t := &Telemetry{
		config:  *config,
		logger:  logger,
		drivers: []Driver{},
	}

	return t, nil
}

// AddDriver adds a new driver to the Telemetry instance.
// Drivers are used to output log messages to different destinations, such as files or network sockets.
// This method is thread-safe.
func (t *Telemetry) AddDriver(driver Driver) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.drivers = append(t.drivers, driver)
	t.logger.(*DefaultLogger).setDrivers(t.drivers)
}

// RemoveDriver removes a driver from the Telemetry instance.
// This method is thread-safe.
func (t *Telemetry) RemoveDriver(driverToRemove Driver) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, driver := range t.drivers {
		if driver == driverToRemove {
			t.drivers = append(t.drivers[:i], t.drivers[i+1:]...)
			break
		}
	}
	t.logger.(*DefaultLogger).setDrivers(t.drivers)
}

// SetLogger sets a new logger for the Telemetry instance.
// The logger is used to output log messages.
func (t *Telemetry) SetLogger(logger Logger) {
	t.logger = logger
}

// Info logs an info message with the provided tags.
// Tags are key-value pairs that provide additional context for the log message.
// Example:
//
//	telemetry.Info("User logged in", map[string]string{"username": "john"})
func (t *Telemetry) Info(message string, tags map[string]string) {
	t.logger.Info(message, tags)
}

// Debug logs a debug message with the provided tags.
// Debug messages are typically used for debugging purposes and are not output by default.
// Example:
//
//	telemetry.Debug("User logged in", map[string]string{"username": "john"})
func (t *Telemetry) Debug(message string, tags map[string]string) {
	t.logger.Debug(message, tags)
}

// Warning logs a warning message with the provided tags.
// Warning messages indicate potential problems or unexpected events.
// Example:
//
//	telemetry.Warning("User login failed", map[string]string{"username": "john"})
func (t *Telemetry) Warning(message string, tags map[string]string) {
	t.logger.Warning(message, tags)
}

// Error logs an error message with the provided tags.
// Error messages indicate serious problems or errors.
// Example:
//
//	telemetry.Error("User login failed", map[string]string{"username": "john"})
func (t *Telemetry) Error(message string, tags map[string]string) {
	t.logger.Error(message, tags)
}

// StartTransaction starts a new transaction with the provided ID and attributes.
// Transactions are used to group related log messages together.
// Example:
//
//	telemetry.StartTransaction("user_login", map[string]string{"username": "john"})
func (t *Telemetry) StartTransaction(transactionID string, attributes map[string]string) {
	t.logger.startTransaction(transactionID, attributes)
}

// EndTransaction ends the transaction with the provided ID.
// Example:
//
//	telemetry.EndTransaction("user_login")
func (t *Telemetry) EndTransaction(transactionID string) {
	t.logger.endTransaction(transactionID)
}
