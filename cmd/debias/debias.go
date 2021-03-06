package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dreadl0ck/debias"
	humanize "github.com/dustin/go-humanize"
)

var (
	flagPath = flag.String("p", "samples-rtl", "file system path with data files")
	flagKaminsky = flag.Bool("k", false, "use kaminsky debiasing")
	flagExt = flag.String("ext", ".wav", "file extension for files to process")
)

func main() {

	flag.Parse()

	var (
		mode debias.Mode
		statsFile string
	)
	if *flagKaminsky {
		mode = debias.ModeKaminsky
		statsFile = filepath.Dir(*flagPath) + "/stats-" + filepath.Base(*flagPath) + "-kaminsky.csv"
	} else {
		mode = debias.ModeVonNeumann
		statsFile = filepath.Dir(*flagPath) + "/stats-" + filepath.Base(*flagPath) + "-neumann.csv"
	}

	stats := debias.Directory(*flagPath, *flagExt, mode)

	f, err := os.Create(statsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var w = csv.NewWriter(f)

	err = w.Write([]string{
		"fileName",
		"bytesIn",
		"bytesOut",
		"duration",
		"sizeDecrease",
		"inputBytesPerSecond",
		"outputBytesPerSecond",
	})	
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range stats {
		err = w.Write([]string{
			s.FileName,
			strconv.FormatInt(s.BytesIn, 10),
			strconv.FormatInt(s.BytesOut, 10),
			s.Duration.String(),
			
			// size decrease in percent
			"-" + strconv.FormatFloat((1.0 - (float64(s.BytesOut)/float64(s.BytesIn))) * 100, 'f', 2, 64) + "%",
			
			// input bytes per second
			humanize.Bytes(uint64(float64(s.BytesIn) / (float64(s.Duration.Milliseconds() / 1000.0)))) + "/s",

			// input bytes per second
			humanize.Bytes(uint64(float64(s.BytesOut) / (float64(s.Duration.Milliseconds() / 1000.0)))) + "/s",
		})	
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Flush()
}