#!/bin/bash

# Script to capture samples from a configurable frequency band using the HackRF.

# 1000000 = 1 MHz
# 1000000000 = 1 GHz
# 6000000000 = 6 GHz

# - 95000000: 95 MHz (FM)
# - 145800000: 145.80 MHz (ISS)
# - 433000000: 433 MHz (Amateur)
# - 790000000: 790 MHz (MFCN, PPDR)
# - 862000000: 862 MHz
# - 1300000000: 1300 MHz (Aeronautical)
# - 1559000000: 1559 MHz (GNSS: Glonass, Galileo)
# - 2200000000: 2200 MHz (Space research, Radio Astronomy)


declare -a arr=(
    "1000000"
    "5000000"
    "10000000"
    "74000000"
    "95000000"
    "145800000"
    "433000000"
    "790000000"
    "862000000"
    "1300000000"
    "1559000000"
    "2200000000"
)

mkdir -p samples-rtl

for f in "${arr[@]}"; do
    echo "capturing frequency $f"
    time rtl_sdr -f $f -s 2500000 -n 500000000 "samples-rtl2/data-$f.wav"
done