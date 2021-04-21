package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"

	rtl "github.com/jpoirier/gortlsdr"

	"log"
	"sync"
)

var logger *zap.SugaredLogger

func main() {

	println(`              
         ( ( o ) ) 
            /   
           /___                    _________
 ________ / __/ ___________ _____________  /
 /  ___/_/ /_ _/  ___/  __  /_  __ \  __  / 
/  /  __  __/    /   / /_/ /_  / / / /_/ /  
__/    /_/   /__/    \__,_/ /_/ /_/\__,_/   
	`)

	flag.Parse()

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger = l.Sugar()

	//---------- Device Check ----------
	if c := rtl.GetDeviceCount(); c == 0 {
		fmt.Println("No devices found, exiting.")
		os.Exit(0)
	} else {
		for i := 0; i < c; i++ {
			m, p, s, err := rtl.GetDeviceUsbStrings(i)
			if err == nil {
				err = errors.New("")
			}
			logger.Debugf("GetDeviceUsbStrings %s - %s %s %s\n", err, m, p, s)
		}
	}
	indexID := 0
	logger.Debugf("===== Device name: %s =====", rtl.GetDeviceName(indexID))
	logger.Debugf("===== Running tests using device indx: 0 =====")

	uatDev := &UAT{}
	if err := uatDev.sdrConfig(indexID); err != nil {
		log.Fatalf("uatDev = &UAT{indexID: id} failed: %s", err.Error())
	}
	uatDev.wg = &sync.WaitGroup{}
	uatDev.wg.Add(1)
	logger.Debugf("\n======= CTRL+C to exit... =======\n\n")
	go uatDev.read()
	uatDev.sigAbort()
}
