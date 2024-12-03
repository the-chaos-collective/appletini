package main

import "github.com/charmbracelet/log"

type Globals struct {
	ConfigPath string
	LogLevel   log.Level
}

func globals() Globals {
	return Globals{
		ConfigPath: "config.json",
		LogLevel:   log.DebugLevel,
	}
}
