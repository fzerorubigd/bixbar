package main

import (
	"fmt"

	"io"

	"github.com/fzerorubigd/bixbar"
	"github.com/fzerorubigd/bixbar/plugins/psutil_plugin/internal/ipblock"
)

type plugin struct {
}

func (pl *plugin) Initialize(io.Writer) {
}

func (pl plugin) Name() string {
	return "psutil"
}

func (pl plugin) Blocks() []string {
	return []string{"ip"}
}

func (pl plugin) Instance(name string, ins string, cfg map[string]interface{}) (bixbar.SimpleBlock, error) {
	switch name {
	case "ip":
		iface := getString(cfg, "interface", "")
		tpl := getString(cfg, "template", "")
		return ipblock.NewIPBlock(iface, tpl), nil
	default:
		return nil, fmt.Errorf("the block name is invaid : %s", name)
	}
}

// NewBixbarPlugin is the plugin entry point
func NewBixbarPlugin() bixbar.Plugin {
	return &plugin{}
}
