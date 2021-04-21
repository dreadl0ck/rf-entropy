package main

import (
	"errors"
	"flag"
	"fmt"

	"go.uber.org/zap"

	rtl "github.com/jpoirier/gortlsdr"

	"log"
	"sync"
)

var logger *zap.SugaredLogger

func main() {

	flag.Parse()

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger = l.Sugar()

	//---------- Device Check ----------
	if c := rtl.GetDeviceCount(); c == 0 {
		logger.Fatal("No devices found, exiting.")
	} else {
		for i := 0; i < c; i++ {
			m, p, s, err := rtl.GetDeviceUsbStrings(i)
			if err == nil {
				err = errors.New("")
			}
			logger.Infof("GetDeviceUsbStrings %s - %s %s %s\n", err, m, p, s)
		}
	}
	indexID := 0
	logger.Infof("===== Device name: %s =====", rtl.GetDeviceName(indexID))
	logger.Infof("===== Running tests using device indx: 0 =====")

	uatDev := &UAT{}
	if err := uatDev.sdrConfig(indexID); err != nil {
		log.Fatalf("uatDev = &UAT{indexID: id} failed: %s", err.Error())
	}
	uatDev.wg = &sync.WaitGroup{}
	uatDev.wg.Add(1)
	fmt.Printf("\n======= CTRL+C to exit... =======\n\n")
	go uatDev.read()
	uatDev.sigAbort()
}
