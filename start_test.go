package main

import (
	"strings"
	"testing"
)

func TestParseCLI(t *testing.T) {
	for _, tc := range [][]string{
		{"andy"},
		{"andy", "config", "-version"},
		{"andy", "config", "init"},
		{"andy", "config", "generate"},
		{"andy", "config", "validate"},
	} {
		name := strings.Join(tc, " ")
		t.Run(name, func(t *testing.T) {
			if _, err := ParseCLICommands(); err != nil {
				t.Error(err)
			}
		})
	}
}
