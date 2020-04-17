#!/usr/bin/env bash

# tweak profile
echo 'alias ll="ls -lahF"' >> ~/.profile
echo 'export PS1="\$ "' >> ~/.profile
echo 'PATH=$PATH:/opt/gopath/src/golang-bpf/bin' >> ~/.profile

echo <<EOF > ~/.gdbinit
set disassembly-flavor intel
set disassemble-next-line on
add-auto-load-safe-path /usr/local/go/src/runtime/runtime-gdb.py
EOF
