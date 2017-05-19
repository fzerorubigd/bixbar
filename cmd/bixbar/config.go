package main

import (
	"path/filepath"

	"fmt"

	"github.com/fzerorubigd/expand"
	"gopkg.in/fzerorubigd/onion.v3"
)

func loadConfigs() (*onion.Onion, error) {
	o := onion.New()

	home, err := expand.HomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ".bixbar/*.*")

	pwd, err := expand.Pwd()
	if err != nil {
		return nil, err
	}

	pluginDir := []string{
		filepath.Join(home, ".bixbar", "plugins"),
		filepath.Join(pwd, "plugins"),
	}

	def := onion.NewDefaultLayer()
	def.SetDefault("plugin.folder", pluginDir)
	_ = o.AddLayer(def)

	files, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}
	for i := range files {
		// Try to load each file
		fl := onion.NewFileLayer(files[i])
		// err is not important here
		_ = o.AddLayer(fl)
	}

	return o, nil
}

func castToMap(in interface{}) (map[string]interface{}, error) {
	if m, ok := in.(map[string]interface{}); ok {
		return m, nil
	}

	mi, ok := in.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("the key is not a map")
	}
	l := make(map[string]interface{})
	for i := range mi {
		if s, ok := i.(string); ok {
			l[s] = mi[i]
		}
	}

	return l, nil
}

func getBlocks(cfg *onion.Onion) (map[string]map[string]interface{}, error) {
	blocks, ok := cfg.Get("blocks")
	if !ok {
		return nil, fmt.Errorf("the blocks section is not available inside configs")
	}

	m, err := castToMap(blocks)
	if err != nil {
		return nil, err
	}

	res := make(map[string]map[string]interface{})
	for i := range m {
		cs, err := castToMap(m[i])
		if err == nil {
			res[i] = cs
		}
	}

	return res, nil
}
