package main

import (
	"errors"
	"strings"
)

type Mode string

var (
	UnknownMode    Mode = "unknown"
	PluginMode     Mode = "plugin"
	StandaloneMode Mode = "standalone"
	WebhookMode    Mode = "webhook"

	stringModeMap = map[string]Mode{
		"plugin":     PluginMode,
		"standalone": StandaloneMode,
		"webhook":    WebhookMode,
	}
)

func ParseWorkMode(mode string) (Mode, error) {
	mode = strings.ToLower(mode)
	m, ok := stringModeMap[mode]
	if ok {
		return m, nil
	}
	return UnknownMode, errors.New("Unknown mode")
}
