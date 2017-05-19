package main

import (
	"plugin"

	"path/filepath"

	"fmt"

	"github.com/fzerorubigd/bixbar"
	"gopkg.in/fzerorubigd/onion.v3"
)

type single struct {
	file string
	pl   *plugin.Plugin
	sym  plugin.Symbol

	fn func() bixbar.Plugin

	obj bixbar.Plugin
}

type pluginManager struct {
	config *onion.Onion

	plugins map[string]single

	blocks map[string][]string
}

func (pm *pluginManager) createInstance(ins string, cfg map[string]interface{}) (bixbar.SimpleBlock, string, string, error) {
	name, ok := cfg["name"].(string)
	if !ok {
		return nil, "", "", fmt.Errorf("name is not available in blocks.%s", ins)
	}

	plgnName, ok := cfg["plugin"].(string)
	if !ok {
		plgnName = name
	}

	plug, ok := pm.plugins[plgnName]
	if !ok {
		return nil, "", "", fmt.Errorf("no plugin with name %s is available", plgnName)
	}

	var valid bool
	for _, i := range pm.blocks[plgnName] {
		if i == name {
			valid = true
			break
		}
	}
	if !valid {
		return nil, "", "", fmt.Errorf("plugin %s dose not have any block with name %s", plgnName, name)
	}

	b, err := plug.obj.Instance(name, ins, cfg)
	if err != nil {
		return nil, "", "", err
	}

	return b, name, ins, nil
}

func loadPlugins(cfg *onion.Onion, dir ...string) (*pluginManager, error) {
	pm := &pluginManager{
		config:  cfg,
		plugins: make(map[string]single),
		blocks:  make(map[string][]string),
	}
	for d := range dir {
		so := filepath.Join(dir[d], "*.so")
		files, err := filepath.Glob(so)
		if err != nil {
			return nil, err
		}

		for i := range files {
			var err error
			s := single{
				file: files[i],
			}
			s.pl, err = plugin.Open(s.file)
			if err != nil {
				return nil, err
			}

			s.sym, err = s.pl.Lookup("NewBixbarPlugin")
			if err != nil {
				return nil, err
			}
			var ok bool
			s.fn, ok = s.sym.(func() bixbar.Plugin)
			if !ok {
				return nil, fmt.Errorf("the function sign is not correct for %s", s.file)
			}

			s.obj = s.fn()
			s.obj.Initialize()

			n := s.obj.Name()
			pm.plugins[n] = s
			pm.blocks[n] = s.obj.Blocks()
		}
	}

	return pm, nil
}
