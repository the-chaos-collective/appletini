package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

const NAME string = "Appletini"

var envPrefix string

func init() {
	envPrefix = toUpperSnake(NAME) + "_"
}

type Globals struct {
	ConfigPath string
	LogLevel   log.Level
	LogPrefix  string
}

func globals() Globals {
	return Globals{
		ConfigPath: "config.work.json",
		LogLevel:   strToLogLevel(env("LOG_LEVEL")),
		LogPrefix:  NAME,
	}
}

func env(key string) string {
	return os.Getenv(envPrefix + key)
}

func strToLogLevel(str string) log.Level {
	level, err := log.ParseLevel(str)
	if err != nil {
		return log.InfoLevel
	}
	return level
}

func toUpperSnake(in string) string {
	out := strings.Trim(in, " \n\r")
	out = strings.ReplaceAll(out, " ", "_")
	out = strings.ToUpper(out)
	return out
}
