package logging_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/log"

	"appletini/logging"
)

func TestCharm(t *testing.T) {
	testCases := []struct {
		name     string
		log      func(logging.Logger)
		expected string
	}{
		{name: "debug", log: func(l logging.Logger) { l.Debug("d-msg", "k", "v") }, expected: "d-msg"},
		{name: "info", log: func(l logging.Logger) { l.Info("i-msg") }, expected: "i-msg"},
		{name: "warn", log: func(l logging.Logger) { l.Warn("w-msg") }, expected: "w-msg"},
		{name: "error", log: func(l logging.Logger) { l.Error("e-msg", "err", "boom") }, expected: "boom"},
		{name: "print", log: func(l logging.Logger) { l.Print("p-msg") }, expected: "p-msg"},
		{name: "debugf", log: func(l logging.Logger) { l.Debugf("d-%d", 1) }, expected: "d-1"},
		{name: "infof", log: func(l logging.Logger) { l.Infof("i-%d", 2) }, expected: "i-2"},
		{name: "warnf", log: func(l logging.Logger) { l.Warnf("w-%d", 3) }, expected: "w-3"},
		{name: "errorf", log: func(l logging.Logger) { l.Errorf("e-%d", 4) }, expected: "e-4"},
		{name: "printf", log: func(l logging.Logger) { l.Printf("p-%d", 5) }, expected: "p-5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			//exhaustruct:ignore
			logger := logging.NewCharm(log.NewWithOptions(&buf, log.Options{Level: log.DebugLevel}))

			tc.log(logger)

			if !strings.Contains(buf.String(), tc.expected) {
				t.Fatalf("output %q does not contain %q", buf.String(), tc.expected)
			}
		})
	}
}

func TestCharmRespectsLevel(t *testing.T) {
	var buf bytes.Buffer

	//exhaustruct:ignore
	logger := logging.NewCharm(log.NewWithOptions(&buf, log.Options{Level: log.WarnLevel}))
	logger.Info("hidden")
	logger.Debug("hidden")

	if buf.Len() != 0 {
		t.Fatalf("expected no output below warn level, got %q", buf.String())
	}
}
