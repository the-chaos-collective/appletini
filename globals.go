package main

type Globals struct {
	ConfigPath string
}

func LoadGlobals() Globals {
	return Globals{
		ConfigPath: "config.json",
	}
}
