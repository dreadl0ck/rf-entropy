package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
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
		statsFile = "stats-" + *flagPath + "-km.csv"
	} else {
		mode = debias.ModeVonNeumann
		statsFile = "stats-" + *flagPath + "-vn.csv"
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
			strconv.Itoa(s.BytesIn),
			strconv.Itoa(s.BytesOut),
			s.Duration.String(),
			
			// size decrease in percent
			"-" + strconv.FormatFloat((1.0 - (float64(s.BytesOut)/float64(s.BytesIn))) * 100, 'f', 2, 64) + "%",
			
			// input bytes per second
			humanize.Bytes(uint64(s.BytesIn / int(s.Duration.Milliseconds() / 1000))) + "/s",

			// input bytes per second
			humanize.Bytes(uint64(s.BytesOut / int(s.Duration.Milliseconds() / 1000))) + "/s",
		})	
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Flush()
}