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

## Generate Entropy

https://pthree.org/2015/06/16/hardware-rng-through-an-rtl-sdr-dongle/

## validate

https://en.wikipedia.org/wiki/FIPS_140-2