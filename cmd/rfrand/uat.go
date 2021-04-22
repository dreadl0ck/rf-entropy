package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/dreadl0ck/debias"
	"github.com/dustin/go-humanize"
	rtl "github.com/jpoirier/gortlsdr"
)

var (
	dumpFile *os.File
	rawFile *os.File
	errDumpFile error
	errRawFile error
)

// UAT holds a device context.
type UAT struct {
	dev *rtl.Context
	wg  *sync.WaitGroup
}

// read does synchronous specific reads.
func (u *UAT) read() {
	defer u.wg.Done()
	logger.Debug("Entered UAT read()")

	var (
		buf         bytes.Buffer
		readCnt     uint64
		buffer      = make([]uint8, rtl.DefaultBufLength)
		out         *io.PipeReader

		start = time.Now()
	)

	debias.MaxChunkSize = *flagMaxChunkSize

	if *flagKaminsky {
		fmt.Println("Using Kaminsky debiasing")
		out, _, _ = debias.Kaminsky(&buf, true, int64(debias.MaxChunkSize))
	} else {
		fmt.Println("Using Von Neumann debiasing")
		out, _, _ = debias.VonNeumann(&buf, true)
	}

	if *flagWriteFile != "" {
		fmt.Println("Writing into file:", *flagWriteFile)
		dumpFile, errDumpFile = os.Create(*flagWriteFile)
		if errDumpFile != nil {
			log.Fatal(errDumpFile)
		}
	}
	if *flagWriteRaw {
		fmt.Println("Writing into raw file:", *flagWriteFile + ".raw")
		rawFile, errRawFile = os.Create(*flagWriteFile + ".raw")
		if errRawFile != nil {
			log.Fatal(errRawFile)
		}
	}
	fmt.Println("Chunk Size:", humanize.Bytes(uint64(*flagMaxChunkSize)))

	fmt.Println()
	fmt.Println("==================================================")
	fmt.Println("Rate      ", "Total    ", "Entropy    ", "Duration")
	fmt.Println("==================================================")

	go func() {

		var (
			b   = make([]byte, debias.MaxChunkSize)
			err error
			n   int
			numBytes int64
			numBytesTotal  int64
			lastBlockOkay bool
			windowStart = time.Now()
		)
		for {
			n, err = out.Read(b)
			if err != nil {
				log.Println(err)
				return
			}

			numBytes += int64(n)
			numBytesTotal += int64(n)

			t := time.Now()
			e := entropy(b[:n])
			outputEntropy = append(outputEntropy, &measurement{
				value: e,
				time: t,
			})
			
			if *flagEntropyGuard != 0 {
				if e < float64(*flagEntropyGuard) {
					if lastBlockOkay {
						fmt.Println("\n[entropy-guard] insufficient entropy detected, discarding data block:", e)
					} else {
						fmt.Println("[entropy-guard] insufficient entropy detected, discarding data block:", e)
					}

					if *flagHexDump {
						fmt.Println(hex.Dump(b[:n]))
					}
					
					lastBlockOkay = false
					continue
				} else {
					lastBlockOkay = true
				}
			}

			if *flagHexDump {
				fmt.Println(hex.Dump(b[:n]))
			}
			
			if *flagWriteFile != "" {
				if *flagMaxFileSize > 0 {
					if int(numBytesTotal) >= *flagMaxFileSize {
						current := numBytesTotal-int64(n)
						free := int64(*flagMaxFileSize) - current
						_, err = dumpFile.Write(b[:free])
						if err != nil {
							log.Fatal(err)
						}
						u.shutdown()
						fmt.Println("done! captured", current+free, "bytes in", time.Since(start))
						os.Exit(0)
					}
				}
				_, err = dumpFile.Write(b[:n])
				if err != nil {
					log.Fatal(err)
				}
			}

			if !*flagInputRate {
				r := float64(numBytes)/(float64(time.Since(windowStart).Milliseconds()) / float64(1000.0))
				outputRates = append(outputRates, &measurement{
					value: r,
					time: t,
				})
				clearLine()
				
				fmt.Print(
					pad(humanize.Bytes(uint64(r))+ "/s   ", 7),
					pad(humanize.Bytes(uint64(numBytesTotal)), 7),
					"   ",
					pad(strconv.FormatFloat(e, 'f', 2, 64), 7),
					"     ",
					time.Since(start),
				)
			}

			// reset start and numBytes for average data rate calculation
			if time.Since(windowStart) > *flagRateInterval {
				windowStart = time.Now()
				numBytes = 0
			}
		}
	}()

	var (
		numBytesRaw int64
		numBytesRawTotal int64
		rawBytesWindowStart = time.Now()
	)

	for {
		nRead, err := u.dev.ReadSync(buffer, rtl.DefaultBufLength)
		if err != nil {
			logger.Debugf("ReadSync Failed - error: %s", err)
			break
		}
		
		// logger.Debugf("ReadSync %d", nRead)
		if nRead > 0 {

			// populate buffer
			buf.Write(buffer[:nRead])

			if *flagWriteRaw {
				_, err = rawFile.Write(buffer[:nRead])
				if err != nil {
					log.Fatal(err)
				}
			}

			numBytesRaw += int64(nRead)
			numBytesRawTotal += int64(nRead)
			
			t := time.Now()
			e := entropy(buffer[:nRead])
			inputEntropy = append(inputEntropy, &measurement{
				value: e,
				time: t,
			})
			
			r := float64(numBytesRaw)/(float64(time.Since(rawBytesWindowStart).Milliseconds()) / float64(1000.0))
			inputRates = append(inputRates, &measurement{
				value: r,
				time: t,
			})
			
			if *flagInputRate {
				clearLine()
				fmt.Print(
					pad(humanize.Bytes(uint64(float64(numBytesRaw)/(float64(time.Since(rawBytesWindowStart).Milliseconds()) / float64(1000.0))))+ "/s   ", 7),
					pad(humanize.Bytes(uint64(numBytesRawTotal)), 7),
					"   ",
					pad(strconv.FormatFloat(e, 'f', 2, 64), 7),
					"     ",
					time.Since(start),
				)
	
				// reset start and numBytes for average data rate calculation
				if time.Since(rawBytesWindowStart) > *flagRateInterval {
					rawBytesWindowStart = time.Now()
					numBytesRaw = 0
				}
			}

			//fmt.Printf("\rnRead %d: readCnt: %d", nRead, readCnt)
			readCnt++
		}
	}
}

// shutdown
func (u *UAT) shutdown() {
	logger.Debug("\nEntered UAT shutdown() ...")
	logger.Debug("UAT shutdown(): closing device ...")
	logger.Debug("u.dev.Close():", u.dev.Close()) // preempt the blocking ReadSync call
	logger.Debug("UAT shutdown(): calling uatWG.Wait() ...")
	u.wg.Wait() // Wait for the goroutine to shutdown
	logger.Debug("UAT shutdown(): uatWG.Wait() returned...")
	fmt.Println() // add newline
	// close file handles
	if dumpFile != nil {
		errDumpFile = dumpFile.Close()
		if errDumpFile != nil {
			log.Fatal(errDumpFile)
		}
	}
	if rawFile != nil {
		errRawFile = rawFile.Close()
		if errRawFile != nil {
			log.Fatal(errRawFile)
		}
	}

	makePlot("input-rates.png", inputRates, "time", "rate", true)
	makePlot("output-rates.png", outputRates, "time", "rate", true)
	makePlot("input-entropy.png", inputEntropy, "time", "entropy", false)
	makePlot("output-entropy.png", outputEntropy, "time", "entropy", false)
}

// sdrConfig configures the device to 978 MHz UAT channel.
func (u *UAT) sdrConfig(indexID int) (err error) {
	if u.dev, err = rtl.Open(indexID); err != nil {
		logger.Debugf("UAT Open Failed...")
		return
	}
	logger.Debugf("GetTunerType: %s", u.dev.GetTunerType())

	// ---------- Set Tuner Gain ----------
	err = u.dev.SetTunerGainMode(true)
	if err != nil {
		u.dev.Close()
		logger.Debugf("SetTunerGainMode Failed - error: %s", err)
		return
	}
	logger.Debugf("SetTunerGainMode Successful")

	var tgain = 0
	gains, err := u.dev.GetTunerGains()
	if err != nil {
		logger.Debugf("GetTunerGains Failed - error: %s", err)
	} else if len(gains) > 0 {
		tgain = int(gains[0])
	}

	// allow gain overwrite
	if *flagTunerGain != 0 {
		tgain = *flagTunerGain
	}

	logger.Debugf("Using gain: %s", tgain)

	err = u.dev.SetTunerGain(tgain)
	if err != nil {
		u.dev.Close()
		logger.Debugf("SetTunerGain Failed - error: %s", err)
		return
	}
	logger.Debugf("SetTunerGain Successful")

	// ---------- Get/Set Sample Rate ----------
	err = u.dev.SetSampleRate(*flagSampleRate)
	if err != nil {
		u.dev.Close()
		logger.Debugf("SetSampleRate Failed - error: %s", err)
		return
	}
	logger.Debug("SetSampleRate - rate: %d", *flagSampleRate)
	logger.Debugf("GetSampleRate: %d", u.dev.GetSampleRate())

	// ---------- Get/Set Xtal Freq ----------
	rtlFreq, tunerFreq, err := u.dev.GetXtalFreq()
	if err != nil {
		u.dev.Close()
		logger.Debugf("GetXtalFreq Failed - error: %s", err)
		return
	}
	logger.Debugf("GetXtalFreq - Rtl: %d, Tuner: %d", rtlFreq, tunerFreq)

	err = u.dev.SetXtalFreq(*flagRTLFreq, *flagTunerFreq)
	if err != nil {
		u.dev.Close()
		logger.Debugf("SetXtalFreq Failed - error: %s", err)
		return
	}
	logger.Debugf("SetXtalFreq - Center freq: %d, Tuner freq: %d",
		*flagRTLFreq, *flagTunerFreq)

	// ---------- Get/Set Center Freq ----------
	err = u.dev.SetCenterFreq(*flagFrequency)
	if err != nil {
		u.dev.Close()
		logger.Debugf("SetCenterFreq Failed, error: %s", err)
		return
	}
	logger.Debugf("SetCenterFreq Successful")

	logger.Debugf("GetCenterFreq: %d", u.dev.GetCenterFreq())
	
	fmt.Println("Frequency: ", *flagFrequency, "Hz")

	// ---------- Set Bandwidth ----------
	logger.Debugf("Setting Bandwidth: %d", *flagBandwidth)
	if err = u.dev.SetTunerBw(*flagBandwidth); err != nil {
		u.dev.Close()
		logger.Debugf("SetTunerBw %d Failed, error: %s", *flagBandwidth, err)
		return
	}
	logger.Debugf("SetTunerBw %d Successful", *flagBandwidth)

	if err = u.dev.ResetBuffer(); err != nil {
		u.dev.Close()
		logger.Debugf("ResetBuffer Failed - error: %s", err)
		return
	}
	logger.Debugf("ResetBuffer Successful")

	// ---------- Get/Set Freq Correction ----------
	freqCorr := u.dev.GetFreqCorrection()
	logger.Debugf("GetFreqCorrection: %d", freqCorr)
	err = u.dev.SetFreqCorrection(freqCorr)
	if err != nil {
		u.dev.Close()
		logger.Debugf("SetFreqCorrection %d Failed, error: %s", freqCorr, err)
		return
	}
	logger.Debugf("SetFreqCorrection %d Successful", freqCorr)

	return
}

// sigAbort
func (u *UAT) sigAbort() {
	
	// setup signal channel
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	
	u.shutdown()
	os.Exit(0)
}

