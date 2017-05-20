package bixbar

import "io"

// Plugin is the block plugin
type Plugin interface {
	// Initialize the plugin. called right after loading th plugin, the writer is fro log.
	Initialize(io.Writer)
	// Name is the plugin name, must be unique
	Name() string
	// Blocks return the available block for this plugin, each name in array
	// should be unique in one plugin
	Blocks() []string
	// Instance create a new block, the first parameter is one of the block names return from
	// blocks function, the first parameter is its name and the next is instance
	Instance(string, string, map[string]interface{}) (SimpleBlock, error)
}
