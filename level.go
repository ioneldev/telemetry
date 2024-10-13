package telemetry

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warning
	Error
)

func (l LogLevel) String() string {
	switch l {
	case Debug:
		return "Debug"
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	default:
		return "Unknown"
	}
}
