#!/bin/bash

search_dir=./maps/solvables
for entry in "$search_dir"/*
do
    echo "test : $entry"
    go run . -f "$entry"
done