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

