#!/bin/bash

for f in samples-rtl/*; do
    echo "processing $f"
    dd if="$f" of="$f-125MB.wav" bs=1 count=125000000 &
done

for f in samples/*; do
    echo "processing $f"
    dd if="$f" of="$f-125MB.wav" bs=1 count=125000000 &
done