package logging

import (
	"github.com/charmbracelet/log"
)

type Charm struct {
	logger *log.Logger
}

func NewCharm(logger *log.Logger) Logger {
	return &Charm{
		logger: logger,
	}
}

func (c *Charm) Debug(msg interface{}, keyvals ...interface{}) {
	c.logger.Debug(msg, keyvals...)
}

// Info prints an info message.
func (c *Charm) Info(msg interface{}, keyvals ...interface{}) {
	c.logger.Info(msg, keyvals...)
}

// Warn prints a warning message.
func (c *Charm) Warn(msg interface{}, keyvals ...interface{}) {
	c.logger.Warn(msg, keyvals...)
}

// Error prints an error message.
func (c *Charm) Error(msg interface{}, keyvals ...interface{}) {
	c.logger.Error(msg, keyvals...)
}

// Fatal prints a fatal message and exits.
func (c *Charm) Fatal(msg interface{}, keyvals ...interface{}) {
	c.logger.Fatal(msg, keyvals...)
}

// Print prints a message with no level.
func (c *Charm) Print(msg interface{}, keyvals ...interface{}) {
	c.logger.Print(msg, keyvals...)
}

// Debugf prints a debug message with formatting.
func (c *Charm) Debugf(format string, args ...interface{}) {
	c.logger.Debugf(format, args...)
}

// Infof prints an info message with formatting.
func (c *Charm) Infof(format string, args ...interface{}) {
	c.logger.Infof(format, args...)
}

// Warnf prints a warning message with formatting.
func (c *Charm) Warnf(format string, args ...interface{}) {
	c.logger.Warnf(format, args...)
}

// Errorf prints an error message with formatting.
func (c *Charm) Errorf(format string, args ...interface{}) {
	c.logger.Errorf(format, args...)
}

// Fatalf prints a fatal message with formatting and exits.
func (c *Charm) Fatalf(format string, args ...interface{}) {
	c.logger.Fatalf(format, args...)
}

// Printf prints a message with no level and formatting.
func (c *Charm) Printf(format string, args ...interface{}) {
	c.logger.Printf(format, args...)
}
