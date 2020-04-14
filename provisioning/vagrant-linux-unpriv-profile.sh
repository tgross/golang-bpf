#!/usr/bin/env bash

# tweak profile
echo 'alias ll="ls -lahF"' >> ~/.profile
echo 'export PS1="\$ "' >> ~/.profile

echo "set disassembly-flavor intel" > ~/.gdbinit
echo "add-auto-load-safe-path /usr/local/go/src/runtime/runtime-gdb.py" >> ~/.gdbinit
