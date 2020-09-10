package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWorkMode(t *testing.T) {
	tests := []struct {
		in      string
		out     Mode
		wantErr bool
	}{
		{
			wantErr: true,
			out:     UnknownMode,
		},
		{
			in:  "plugin",
			out: PluginMode,
		},
	}
	for _, test := range tests {
		out, err := ParseWorkMode(test.in)
		assert.Equal(t, test.wantErr, err != nil)
		assert.Equal(t, test.out, out)
	}
}
