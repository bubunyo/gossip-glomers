#!/bin/bash
# set -e

echo Running maelstrom for "$1" as "$2" ...
cd $1
echo Remove old binary...
rm -f ./$2
echo Compiling binary
go mod tidy && go build -o $2
echo Binary compile done
cd ..
p=$1
t=$2
nc=$3
r=$4
shift 4
./maelstrom/maelstrom test -w $t \
  --bin ./$p/$t   \
  --time-limit 30 \
  --node-count $nc \
  --rate $r ${@:+"$@"}
  
