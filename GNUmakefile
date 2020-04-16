MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -euc
.DEFAULT_GOAL := help

BUILD_DIR := $(shell pwd)/bin

.PHONY: clean distclean helloworld minimal web nomad

## display this help message
help:
	@echo -e "\033[32m"
	@echo "Targets in this Makefile set up test targets."
	@echo
	@awk '/^##.*$$/,/[a-zA-Z_-]+:/' $(MAKEFILE_LIST) | awk '!(NR%2){print $$0p}{p=$$0}' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

## build the helloworld binary
helloworld: bin/helloworld
bin/helloworld: targets/helloworld/*
	mkdir -p ./bin
	cd targets/helloworld && go build -o $(BUILD_DIR)/helloworld

## build the minimal binary
minimal: bin/minimal
bin/minimal: targets/minimal/*
	mkdir -p ./bin
	cd targets/minimal && go build -gcflags '-l' -o $(BUILD_DIR)/minimal

## build the worker binary
worker: bin/worker
bin/worker: targets/worker/*
	mkdir -p ./bin
	cd targets/worker && go build -gcflags '-l' -o $(BUILD_DIR)/worker

## build the web binary
web: bin/web
bin/web: targets/web/*
	mkdir -p ./bin
	cd targets/web && go build -o $(BUILD_DIR)/web

## download Nomad
nomad: bin/nomad

bin/nomad: targets/nomad/nomad.zip
	mkdir -p ./bin
	unzip "targets/nomad/nomad.zip" -d targets/nomad
	mv ./targets/nomad/nomad ./bin/nomad

targets/nomad/nomad.zip:
	mkdir -p ./targets/nomad
	curl -o "targets/nomad/nomad.zip" \
		"https://releases.hashicorp.com/nomad/0.11.0/nomad_0.11.0_linux_amd64.zip"

## clean up the bin directory
clean:
	rm -rf ./bin

## clean up the bin directory and downloaded files
distclean: clean
	rm -rf ./targets/nomad
