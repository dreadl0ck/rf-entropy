#!/bin/bash

# Script to capture samples from a configurable frequency band using the HackRF.

# 1000000 = 1 MHz
# 1000000000 = 1 GHz
# 6000000000 = 6 GHz

sleep=55

# - 95000000: 95 MHz (FM)
# - 145800000: 145.80 MHz (ISS)
# - 433000000: 433 MHz (Amateur)
# - 790000000: 790 MHz (MFCN, PPDR)
# - 862000000: 862 MHz
# - 1300000000: 1300 MHz (Aeronautical)
# - 1559000000: 1559 MHz (GNSS: Glonass, Galileo)
# - 2200000000: 2200 MHz (Space research, Radio Astronomy)

declare -a arr=("95000000"
        "145800000"
        "433000000"
        "790000000"
        "862000000"
        "1300000000"
        "1559000000"
        "2200000000"
        )

echo "current sleep duration = $sleep"

mkdir -p samples

for f in "${arr[@]}"; do
    echo "capturing frequency $f"
    hackrf_transfer -r samples/data-$f.wav -f $f &

    sleep $sleep

    # kill via signal after sleep duration
    kill -INT $(pgrep hackrf_transfer)

    # give hackrf some time to close handles
    sleep 1
done