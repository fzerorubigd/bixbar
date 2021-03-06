package main

import (
	"fmt"

	"io"

	"github.com/fzerorubigd/bixbar"
)

type plugin struct {
}

func (pl *plugin) Initialize(io.Writer) {
}

func (pl plugin) Name() string {
	return "static"
}

func (pl plugin) Blocks() []string {
	return []string{"static"}
}

func (pl plugin) Instance(name string, ins string, cfg map[string]interface{}) (bixbar.SimpleBlock, error) {
	switch name {
	case "static":
		txt, _ := cfg["text"].(string)
		clrTxt, _ := cfg["color"].(string)
		clr, _ := bixbar.NewColor(clrTxt)
		return &staticBlock{
			fullText: txt,
			color:    clr,
		}, nil
	default:
		return nil, fmt.Errorf("the block name is invaid : %s", name)
	}
}

// NewBixbarPlugin is the plugin entry point
func NewBixbarPlugin() bixbar.Plugin {
	return &plugin{}
}
