#!/bin/bash
RES=0
SIZE=4
ITER=0

while [ $RES -eq 0 ]
do
    python npuzzle-gen.py -s $SIZE  > test.txt && go run . -w 8 -ss 32 -f test.txt
    if [ $ITER -eq 20 ]
    then
        SIZE=$((SIZE+1))
        ITER=0
    else
        ITER=$((ITER+1))
    fi
    RES=$?
    if [ $SIZE -eq 5 ]
    then
        break
    fi
done