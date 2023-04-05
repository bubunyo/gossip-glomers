#!/bin/bash
echo Running maelstrom for "$1" ...
cd $1
go mod tidy && go build
cd ..
./maelstrom/maelstrom test -w echo --bin ./$1/$1 --node-count 1 --time-limit 10
