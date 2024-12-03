package logging

type Logger interface {
	Debug(msg interface{}, keyvals ...interface{})

	// Info prints an info message.
	Info(msg interface{}, keyvals ...interface{})

	// Warn prints a warning message.
	Warn(msg interface{}, keyvals ...interface{})

	// Error prints an error message.
	Error(msg interface{}, keyvals ...interface{})

	// Fatal prints a fatal message and exits.
	Fatal(msg interface{}, keyvals ...interface{})

	// Print prints a message with no level.
	Print(msg interface{}, keyvals ...interface{})

	// Debugf prints a debug message with formatting.
	Debugf(format string, args ...interface{})

	// Infof prints an info message with formatting.
	Infof(format string, args ...interface{})

	// Warnf prints a warning message with formatting.
	Warnf(format string, args ...interface{})

	// Errorf prints an error message with formatting.
	Errorf(format string, args ...interface{})

	// Fatalf prints a fatal message with formatting and exits.
	Fatalf(format string, args ...interface{})

	// Printf prints a message with no level and formatting.
	Printf(format string, args ...interface{})
}
