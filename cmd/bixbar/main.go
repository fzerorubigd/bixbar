package main

import (
	"os/signal"
	"path/filepath"

	"log"

	"os"
	"time"

	"syscall"

	"github.com/fzerorubigd/expand"
	"github.com/ogier/pflag"
	_ "gopkg.in/fzerorubigd/onion.v3/yamlloader"
)

var (
	pd *string
)

func main() {
	pflag.Parse()

	o, err := loadConfigs()
	if err != nil {
		log.Fatal(err)
	}
	dir := o.GetStringSlice("plugin.folder")
	pm, err := loadPlugins(o, dir...)
	if err != nil {
		log.Fatal(err)
	}

	data, err := getBlocks(o)
	if err != nil {
		log.Fatal(err)
	}
	refresh := o.GetDurationDefault("bixbar.refresh", time.Second)
	bar := NewBar(refresh, os.Stdout, os.Stdin)
	for i := range data {
		b, name, ins, err := pm.createInstance(i, data[i])
		if err != nil {
			log.Fatal(err)
		}
		bar.AddBlock(name, ins, b)
	}

	bar.Start()

	quit := make(chan os.Signal, 6)
	signal.Notify(quit, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT)
	<-quit

	bar.Stop()
}

func init() {
	pwd, _ := expand.Pwd()
	pd = pflag.StringP(
		"plugin-folder",
		"p",
		filepath.Join(pwd, "plugins"),
		"the plugin folder. all file with so extension are loaded in this directory",
	)
}
