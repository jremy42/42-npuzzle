#!/bin/bash

search_dir=./maps/wrongMap
for entry in "$search_dir"/*
do
    echo "test : $entry"
    go run . -w 8 -split 32 -f "$entry"
done