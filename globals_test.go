package main

import (
	"testing"

	"github.com/charmbracelet/log"
)

func TestToUpperSnake(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "Appletini", expected: "APPLETINI"},
		{input: "apple tini", expected: "APPLE_TINI"},
		{input: "  padded name \n", expected: "PADDED_NAME"},
		{input: "", expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			if got := toUpperSnake(tc.input); got != tc.expected {
				t.Fatalf("got %q, expected %q", got, tc.expected)
			}
		})
	}
}

func TestStrToLogLevel(t *testing.T) {
	testCases := []struct {
		input    string
		expected log.Level
	}{
		{input: "debug", expected: log.DebugLevel},
		{input: "warn", expected: log.WarnLevel},
		{input: "error", expected: log.ErrorLevel},
		{input: "", expected: log.InfoLevel},
		{input: "nonsense", expected: log.InfoLevel},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			if got := strToLogLevel(tc.input); got != tc.expected {
				t.Fatalf("got %v, expected %v", got, tc.expected)
			}
		})
	}
}

func TestGlobals(t *testing.T) {
	testCases := []struct {
		name     string
		envValue string
		expected log.Level
	}{
		{name: "level from env", envValue: "debug", expected: log.DebugLevel},
		{name: "default level", envValue: "", expected: log.InfoLevel},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("APPLETINI_LOG_LEVEL", tc.envValue)

			g := globals()

			if g.LogLevel != tc.expected || g.ConfigPath != "config.json" || g.LogPrefix != NAME {
				t.Fatalf("unexpected globals %+v", g)
			}
		})
	}
}
