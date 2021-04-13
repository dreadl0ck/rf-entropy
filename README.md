# rf-entropy

RF entropy experiments for random number generation

Install deps:

    apt install gnuradio gr-osmosdr rtl-sdr and gqrx-sdr

Plugin device and stop it being used as a dvbt device:
    
    rmmod dvb usb rtl28xxu rtl2832

If this results in an error, stating that the module is currently loaded and busy:

/etc/modprobe.d/blacklist.conf:

    blacklist dvb_usb_rtl28xxu

then reboot for the changes to take effect.

## Setup rtl_entropy

Repo: https://github.com/pwarren/rtl-entropy

Install deps:

> LibCAP not LibPCAP ;-)

    apt install libcap-dev librtlsdr-dev cmake

Compile 

    mkdir build
    cd build
    cmake ../
    make
    sudo make install

usb error -6 ? Device handle has been claimed by another process:

    pkill rtl_entropy

test

    $ sudo rtl_entropy -b -f 74M
    $ tail -f /run/rtl_entropy.fifo | dd of=/dev/null
    ^C8999+10 records in
    9004+0 records out
    4610048 bytes (4.6 MB) copied, 13.294 s, 347 kB/s

Choose frequency:

    entropy_rtl -b -f 74.8M

Fetch data:

    # tail -f /run/rtl_entropy.fifo | dd of=random.img bs=1 count=512000000 iflag=fullblock
    339486772 bytes (339 MB, 324 MiB) copied, 1170 s, 290 kB/s

    511915211 bytes (512 MB, 488 MiB) copied, 1763 s, 290 kB/s
    512000000+0 records in
    512000000+0 records out
    512000000 bytes (512 MB, 488 MiB) copied, 1763.26 s, 290 kB/s
    
RNG test:
    
    # rngtest < random.img 
    rngtest 5
    Copyright (c) 2004 by Henrique de Moraes Holschuh
    This is free software; see the source for copying conditions.  There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

    rngtest: starting FIPS tests...
    rngtest: entropy source drained
    rngtest: bits received from input: 4096000000
    rngtest: FIPS 140-2 successes: 204646
    rngtest: FIPS 140-2 failures: 153
    rngtest: FIPS 140-2(2001-10-10) Monobit: 23
    rngtest: FIPS 140-2(2001-10-10) Poker: 18
    rngtest: FIPS 140-2(2001-10-10) Runs: 48
    rngtest: FIPS 140-2(2001-10-10) Long run: 64
    rngtest: FIPS 140-2(2001-10-10) Continuous run: 0
    rngtest: input channel speed: (min=229.801; avg=19894.427; max=19073.486)Mibits/s
    rngtest: FIPS tests speed: (min=19.150; avg=181.189; max=186.995)Mibits/s
    rngtest: Program run time: 21768269 microseconds

## Generate Entropy

https://pthree.org/2015/06/16/hardware-rng-through-an-rtl-sdr-dongle/

## validate

https://en.wikipedia.org/wiki/FIPS_140-2
