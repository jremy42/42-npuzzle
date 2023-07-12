#!/bin/bash
RES=0
SIZE=40
ITER=0

while [ $RES -eq 0 ]
do
    python npuzzle-gen.py -s $SIZE -i 1  > test.txt && go run . -w 8 -ss 32 -f test.txt
    if [ $? -ne 0 ]
    then
        echo "test failed"
        exit 1
    fi
    if [ $ITER -eq 1 ]
    then
        SIZE=$((SIZE+10))
        ITER=0
    else
        ITER=$((ITER+10))
    fi
    RES=$?
    if [ $SIZE -eq 1000 ]
    then
        break
    fi
done