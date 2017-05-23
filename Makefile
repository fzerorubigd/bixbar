export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
all: bixbar static i3blocks psutil

bixbar:
	cd $(ROOT)/cmd/bixbar && go build

static:
	cd $(ROOT)/plugins/static_plugin && go build -buildmode=plugin

i3blocks:
	cd $(ROOT)/plugins/i3blocks_plugin && go build -buildmode=plugin

psutil:
	cd $(ROOT)/plugins/psutil_plugin && go build -buildmode=plugin

install: all
	mkdir -p ~/.bixbar/plugins
	cp $(ROOT)/plugins/*/*.so ~/.bixbar/plugins
