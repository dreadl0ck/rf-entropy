package main

import (
	"flag"
	"time"
)

var (
	flagKaminsky  = flag.Bool("k", false, "use kaminsky mode")
	flagFrequency = flag.Int("f", 145800000, "set frequency")
	flagSampleRate = flag.Int("s", 2083334, "set sample rate")
	flagBandwidth = flag.Int("b", 1000000, "bandwidth")
	flagRTLFreq = flag.Int("r", 28800000, "rtl frequency")
	flagTunerFreq = flag.Int("t", 28800000, "tuner frequency")
	flagTunerGain = flag.Int("g", 0, "tuner gain")

	flagHexDump = flag.Bool("hex", false, "hexdump")
	flagWriteFile = flag.String("w", "", "write into file")

	flagMaxChunkSize = flag.Int("c", 1024 * 100, "max chunk size")
	flagRateInterval = flag.Duration("i", 5 * time.Second, "rate interval")
	
	flagEntropyGuard = flag.Int("e", 5, "entropy guard")

	flagMaxFileSize = flag.Int("size", 0, "max file size (default unlimited)")
	flagWriteRawInput = flag.Bool("raw", false, "write raw input data into file")
	flagInputRate = flag.Bool("input-rate", false, "show input data rate instead of output data rate")
)
