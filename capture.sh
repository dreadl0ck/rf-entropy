#!/bin/bash

# Script to capture samples from a configurable frequency band using the HackRF.

# 1000000 = 1 MHz
# 1000000000 = 1 GHz
# 6000000000 = 6 GHz

sleep=5
start=1000000
end=6000000000
step=10000000

echo "current sleep duration = $sleep"

duration=$(((($end-$start)/$step)*$sleep))
printf 'processing will take %dh:%dm:%ds\n' $((duration/3600)) $((duration%3600/60)) $((duration%60))

res=$((($end-$start)/$step))

echo "and create $res files"

mkdir -p samples

for (( COUNTER=$start; COUNTER<=$end; COUNTER+=$step )); do
    echo "capturing frequency $COUNTER"
    hackrf_transfer -r samples/data-$COUNTER.wav -f $COUNTER &

    sleep $sleep

    # kill via signal after sleep duration
    kill -INT $(pgrep hackrf_transfer)
done