package main

type Globals struct {
	ConfigPath string
}

func globals() Globals {
	return Globals{
		ConfigPath: "config.json",
	}
}
