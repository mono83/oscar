package events

// LogEvent contains logging event
type LogEvent struct {
	Level   byte
	Pattern string
}

// LevelString returns string representation of log level
func (l LogEvent) LevelString() string {
	switch l.Level {
	case LogLevelTrace:
		return "trace"
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelError:
		return "error"
	default:
		return "unknown"
	}
}

// List of log levels
const (
	LogLevelTrace byte = 0
	LogLevelDebug byte = 1
	LogLevelInfo  byte = 2
	LogLevelError byte = 3
)
