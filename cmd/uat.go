package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/dreadl0ck/debias"
	"github.com/dustin/go-humanize"
	rtl "github.com/jpoirier/gortlsdr"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	dumpFile *os.File
	errDumpFile error
)

// UAT holds a device context.
type UAT struct {
	dev *rtl.Context
	wg  *sync.WaitGroup
}

// read does synchronous specific reads.
func (u *UAT) read() {
	defer u.wg.Done()
	log.Println("Entered UAT read() ...")

	var (
		buf         bytes.Buffer
		readCnt     uint64
		buffer      = make([]uint8, rtl.DefaultBufLength)
		out         *io.PipeReader

		start = time.Now()
	)

	debias.MaxChunkSize = *flagMaxChunkSize

	if *flagKaminsky {
		out, _, _ = debias.Kaminsky(&buf, true, 1024)
	} else {
		out, _, _ = debias.VonNeumann(&buf, true)
	}

	if *flagWriteFile != "" {
		dumpFile, errDumpFile = os.Create(*flagWriteFile)
		if errDumpFile != nil {
			log.Fatal(errDumpFile)
		}
	}

	fmt.Println("Rate      ", "Total   ", "Entropy", "Duration")

	go func() {

		var (
			b   = make([]byte, debias.MaxChunkSize)
			err error
			n   int
			numBytes int64
		)
		for {
			n, err = out.Read(b)
			if err != nil {
				log.Println(err)
				return
			}

			numBytes += int64(n)

			if *flagHexDump {
				fmt.Println(hex.Dump(b[:n]))
			}
			if *flagWriteFile != "" {
				_, err = dumpFile.Write(b[:n])
				if err != nil {
					log.Fatal(err)
				}
			}

			if numBytes != 0 {
				clearLine()
				fmt.Print(
					humanize.Bytes(uint64(float64(numBytes)/float64(time.Since(start).Milliseconds() / 1000.0)))+ "/s   ",
					humanize.Bytes(uint64(numBytes)),
					"   ",
					debias.ShannonEntropy(b[:n]),
					"    ",
					time.Since(start),
				)
			}
		}
	}()

	for {
		nRead, err := u.dev.ReadSync(buffer, rtl.DefaultBufLength)
		if err != nil {
			logger.Infof("\tReadSync Failed - error: %s\n", err)
			break
		}
		// logger.Infof("\tReadSync %d\n", nRead)
		if nRead > 0 {

			// populate buffer
			buf.Write(buffer[:nRead])

			//fmt.Printf("\rnRead %d: readCnt: %d", nRead, readCnt)
			readCnt++
		}
	}
}

// shutdown
func (u *UAT) shutdown() {
	fmt.Println()
	log.Println("\nEntered UAT shutdown() ...")
	log.Println("UAT shutdown(): closing device ...")
	log.Println("u.dev.Close():", u.dev.Close()) // preempt the blocking ReadSync call
	log.Println("UAT shutdown(): calling uatWG.Wait() ...")
	u.wg.Wait() // Wait for the goroutine to shutdown
	log.Println("UAT shutdown(): uatWG.Wait() returned...")
}

// sdrConfig configures the device to 978 MHz UAT channel.
func (u *UAT) sdrConfig(indexID int) (err error) {
	if u.dev, err = rtl.Open(indexID); err != nil {
		logger.Infof("\tUAT Open Failed...\n")
		return
	}
	logger.Infof("\tGetTunerType: %s\n", u.dev.GetTunerType())

	// ---------- Set Tuner Gain ----------
	err = u.dev.SetTunerGainMode(true)
	if err != nil {
		u.dev.Close()
		logger.Infof("\tSetTunerGainMode Failed - error: %s\n", err)
		return
	}
	logger.Infof("\tSetTunerGainMode Successful\n")

	var tgain = 0
	gains, err := u.dev.GetTunerGains()
	if err != nil {
		logger.Infof("\tGetTunerGains Failed - error: %s\n", err)
	} else if len(gains) > 0 {
		tgain = int(gains[0])
	}

	// allow gain overwrite
	if *flagTunerGain != 0 {
		tgain = *flagTunerGain
	}

	logger.Infof("\tUsing gain: %s\n", tgain)

	err = u.dev.SetTunerGain(tgain)
	if err != nil {
		u.dev.Close()
		logger.Infof("\tSetTunerGain Failed - error: %s\n", err)
		return
	}
	logger.Infof("\tSetTunerGain Successful\n")

	// ---------- Get/Set Sample Rate ----------
	err = u.dev.SetSampleRate(*flagSampleRate)
	if err != nil {
		u.dev.Close()
		logger.Infof("\tSetSampleRate Failed - error: %s\n", err)
		return
	}
	logger.Info("\tSetSampleRate - rate: %d\n", *flagSampleRate)
	logger.Infof("\tGetSampleRate: %d\n", u.dev.GetSampleRate())

	// ---------- Get/Set Xtal Freq ----------
	rtlFreq, tunerFreq, err := u.dev.GetXtalFreq()
	if err != nil {
		u.dev.Close()
		logger.Infof("\tGetXtalFreq Failed - error: %s\n", err)
		return
	}
	logger.Infof("\tGetXtalFreq - Rtl: %d, Tuner: %d\n", rtlFreq, tunerFreq)

	err = u.dev.SetXtalFreq(*flagRTLFreq, *flagTunerFreq)
	if err != nil {
		u.dev.Close()
		logger.Infof("\tSetXtalFreq Failed - error: %s\n", err)
		return
	}
	logger.Infof("\tSetXtalFreq - Center freq: %d, Tuner freq: %d\n",
		*flagRTLFreq, *flagTunerFreq)

	// ---------- Get/Set Center Freq ----------
	err = u.dev.SetCenterFreq(*flagFrequency)
	if err != nil {
		u.dev.Close()
		logger.Infof("\tSetCenterFreq Failed, error: %s\n", err)
		return
	}
	logger.Infof("\tSetCenterFreq Successful\n")

	logger.Infof("\tGetCenterFreq: %d\n", u.dev.GetCenterFreq())

	// ---------- Set Bandwidth ----------
	logger.Infof("\tSetting Bandwidth: %d\n", *flagBandwidth)
	if err = u.dev.SetTunerBw(*flagBandwidth); err != nil {
		u.dev.Close()
		logger.Infof("\tSetTunerBw %d Failed, error: %s\n", *flagBandwidth, err)
		return
	}
	logger.Infof("\tSetTunerBw %d Successful\n", *flagBandwidth)

	if err = u.dev.ResetBuffer(); err != nil {
		u.dev.Close()
		logger.Infof("\tResetBuffer Failed - error: %s\n", err)
		return
	}
	logger.Infof("\tResetBuffer Successful\n")

	// ---------- Get/Set Freq Correction ----------
	freqCorr := u.dev.GetFreqCorrection()
	logger.Infof("\tGetFreqCorrection: %d\n", freqCorr)
	err = u.dev.SetFreqCorrection(freqCorr)
	if err != nil {
		u.dev.Close()
		logger.Infof("\tSetFreqCorrection %d Failed, error: %s\n", freqCorr, err)
		return
	}
	logger.Infof("\tSetFreqCorrection %d Successful\n", freqCorr)

	return
}

// sigAbort
func (u *UAT) sigAbort() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	u.shutdown()
	errDumpFile = dumpFile.Close()
	if errDumpFile != nil {
		log.Fatal(errDumpFile)
	}
	os.Exit(0)
}

func clearLine() {
	print("\033[2K\r")
}
