#!/usr/bin/env bash

# install bpftrace example tools
git clone https://github.com/iovisor/bpftrace.git
mkdir -p /home/vagrant/.local/bin
ln -s $(pwd)/bpftrace/tools/*.bt /home/vagrant/.local/bin
