# rf-entropy

Our receiver:

    officially: 48.25-863.25 Mhz
    tested for: 24 - 1766 MHz

http://www.giga.co.za/ocart/index.php?route=product/product&product_id=52

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

Fetch 512MB of data:

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

DieHarder test suite:

    #dieharder -a < random.img 
    #=============================================================================#
    #            dieharder version 3.31.1 Copyright 2003 Robert G. Brown          #
    #=============================================================================#
       rng_name    |rands/second|   Seed   |
            mt19937|  1.22e+08  | 904524014|
    #=============================================================================#
            test_name   |ntup| tsamples |psamples|  p-value |Assessment
    #=============================================================================#
       diehard_birthdays|   0|       100|     100|0.44405096|  PASSED  
          diehard_operm5|   0|   1000000|     100|0.72574609|  PASSED  
      diehard_rank_32x32|   0|     40000|     100|0.83826473|  PASSED  
        diehard_rank_6x8|   0|    100000|     100|0.75883616|  PASSED  
       diehard_bitstream|   0|   2097152|     100|0.52521521|  PASSED  
            diehard_opso|   0|   2097152|     100|0.54327546|  PASSED  
            diehard_oqso|   0|   2097152|     100|0.63433709|  PASSED  
             diehard_dna|   0|   2097152|     100|0.38572275|  PASSED  
    diehard_count_1s_str|   0|    256000|     100|0.48686044|  PASSED  
    diehard_count_1s_byt|   0|    256000|     100|0.50349216|  PASSED  
     diehard_parking_lot|   0|     12000|     100|0.22547510|  PASSED  
        diehard_2dsphere|   2|      8000|     100|0.86305985|  PASSED  
        diehard_3dsphere|   3|      4000|     100|0.00479805|   WEAK   
         diehard_squeeze|   0|    100000|     100|0.78316946|  PASSED  
            diehard_sums|   0|       100|     100|0.13007386|  PASSED  
            diehard_runs|   0|    100000|     100|0.62865098|  PASSED  
            diehard_runs|   0|    100000|     100|0.13647343|  PASSED  
           diehard_craps|   0|    200000|     100|0.96350555|  PASSED  
           diehard_craps|   0|    200000|     100|0.87759170|  PASSED  
     marsaglia_tsang_gcd|   0|  10000000|     100|0.18253008|  PASSED  
     marsaglia_tsang_gcd|   0|  10000000|     100|0.67098980|  PASSED  
             sts_monobit|   1|    100000|     100|0.51172956|  PASSED  
                sts_runs|   2|    100000|     100|0.21252793|  PASSED 
    
    Preparing to run test 207.  ntuple = 0
    Preparing to run test 208.  ntuple = 0
    Preparing to run test 209.  ntuple = 0

> Note: In a consequtive run on the same data, the WEAK test passed. Why?

## Generate Entropy

https://pthree.org/2015/06/16/hardware-rng-through-an-rtl-sdr-dongle/

## Capture raw data:

http://www.aaronscher.com/wireless_com_SDR/RTL_SDR_AM_spectrum_demod.html

e.g.: rtl_sdr -f 433920000 -g 10 -s 2500000 -n 25000000 random.bin

## Validate

NIST Statistical Test Suite is used for randomness testing. 

Use run_stats.sh to run statistics. Because of the 1000 bitstreams hardcoded in the replacing utilities.c and the supplied streamlength for individual bitstreams in the shell script, the data assessed should be at least 100MB in size.

Reports can be parsed using convert_reports_to_spreadsheet.sh to provide an overview when performing tests on large amounts of input files.

For more information on interpretation of STS results, have a look at the NIST documentation.

## create 1MB file containing 0xFF

    for i in {1..1000}; do printf "\x$(printf %x 255)"; done > temp.bin

    for i in {1..1000}; do cat temp.bin >> filetosend.bin; done

## HackRF send data

    for i in {1..3}; do hackrf_transfer -t filetosend.bin -f 433920000; done

## Run Von Neumann Debiasing

    go run debias.go

# Frequencies

- 95 MHz (FM)
- 145.80 MHz (ISS)
- 433 MHz (Amateur)
- 790 MHz (MFCN, PPDR)
- 862 MHz
- 1300 MHz (Aeronautical)
- 1559 MHz (GNSS: Glonass, Galileo)
- 2200 MHz (Space research, Radio Astronomy)
