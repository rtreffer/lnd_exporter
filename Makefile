.PHONY: default
default: help

.PHONY: build
build: 
	go build

.SILENT:
.PHONY: help
help:
	echo
	echo "  make  <target>   <description>"
	echo "         help      # this help"
	echo "         build     # build the exporter"
	echo