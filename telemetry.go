package telemetry

import (
	"fmt"
	"sync"
)

type Telemetry struct {
	config  Config
	logger  Logger
	drivers []Driver
	mu      sync.RWMutex
}

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

func (t *Telemetry) AddDriver(driver Driver) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.drivers = append(t.drivers, driver)
	t.logger.(*DefaultLogger).setDrivers(t.drivers)
}

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

func (t *Telemetry) SetLogger(logger Logger) {
	t.logger = logger
}

func (t *Telemetry) Info(message string, tags map[string]string) {
	t.logger.Info(message, tags)
}

func (t *Telemetry) Debug(message string, tags map[string]string) {
	t.logger.Debug(message, tags)
}

func (t *Telemetry) Warning(message string, tags map[string]string) {
	t.logger.Warning(message, tags)
}

func (t *Telemetry) Error(message string, tags map[string]string) {
	t.logger.Error(message, tags)
}

func (t *Telemetry) StartTransaction(transactionID string, attributes map[string]string) {
	t.logger.startTransaction(transactionID, attributes)
}

func (t *Telemetry) EndTransaction(transactionID string) {
	t.logger.endTransaction(transactionID)
}
