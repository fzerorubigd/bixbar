package main

import (
	"fmt"

	"strconv"

	"io"

	"github.com/fzerorubigd/bixbar"
)

type plugin struct {
}

func (pl *plugin) Initialize(io.Writer) {
}

func (pl plugin) Name() string {
	return "i3blocks"
}

func (pl plugin) Blocks() []string {
	return []string{"i3blocks"}
}

func (pl plugin) Instance(name string, ins string, cfg map[string]interface{}) (bixbar.SimpleBlock, error) {
	switch name {
	case "i3blocks":
		cmd, _ := cfg["command"].(string)
		interval, _ := cfg["interval"].(int)
		if interval == 0 {
			t, _ := cfg["interval"].(string)
			tI, err := strconv.ParseInt(t, 10, 0)
			if err == nil {
				interval = int(tI)
			}
		}
		label, _ := cfg["label"].(string)
		format, _ := cfg["format"].(string)
		return newShellBlock(name, ins, cmd, interval, label, format), nil

	default:
		return nil, fmt.Errorf("the block name is invaid : %s", name)
	}
}

// NewBixbarPlugin is the plugin entry point
func NewBixbarPlugin() bixbar.Plugin {
	return &plugin{}
}
