#!/bin/bash
set -e

echo Running maelstrom for "$1" as "$2" ...
cd $1
echo Remove old binary...
rm -f ./$2
echo Compiling binary
go mod tidy && go build -o $2
echo Binary compile done
cd ..
./maelstrom/maelstrom test -w $2 --bin ./$1/$2 --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
