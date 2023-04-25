#!/bin/bash
echo Running maelstrom for "$1" ...
cd $1
echo Compiling binary
go mod tidy && go build
echo Binary compile done
cd ..
# ./maelstrom/maelstrom test -w echo --bin ./$1/$1 --node-count 1 --time-limit 10
./maelstrom/maelstrom test -w unique-ids --bin ./$1/$1 --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
