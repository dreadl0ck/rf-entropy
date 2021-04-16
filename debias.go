package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	debug = flag.Bool("d", false, "toggle debug mode")
	flagPath = flag.String("p", "samples-rtl", "path with test data files")
)

func main() {

	flag.Parse()
	
	files, err := ioutil.ReadDir(*flagPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		
		if filepath.Ext(f.Name()) != ".wav" {
			continue
		}
		
		file := filepath.Join(*flagPath, f.Name())
		fmt.Println("processing", file)
		
		debiasFile(file, f)
	}
}

// The algorithm works on pairs of bits, and produces output as follows:
// - If the input is "00" or "11", the input is discarded (no output).
// - If the input is "10", output a "1".
// - If the input is "01", output a "0".
func debiasData(data []byte) bytes.Buffer {
	var (
		buf bytes.Buffer
		ob = byte(0)
		bitcount uint
	)

	for _, b := range data {

		if *debug {
			fmt.Printf("byte: %08b\n", b)
		}
	
		for j:=0; j<8; j+= 2 {

			// fmt.Printf(" ===> b: %08b, j: %d, b >> j  : %08b, (b >> j) & 0x01: %08b \n", b, j, b >> j, (b >> j) & 0x01)
			// fmt.Printf(" ===> b: %08b, j: %d, b >> j+1: %08b, (b >> j+1) & 0x01: %08b \n", b, j, b >> j+1, (b >> j+1) & 0x01)
			
			ch := (b >> (7-j)) & 0x01
			ch2 := (b >> (7-(j+1))) & 0x01
			
			if *debug {
				fmt.Println(ch, ch2)
			}
			
			if (ch != ch2) {
			
				if ch == 1 {
					// store a 1 in our bitbuffer
					ob = setBit(ob, 7-bitcount)

					if *debug {
						fmt.Printf("collecting 1, out byte: %08b\n", ob)
					}
				} // else: leave the buffer alone, it's already 0 at this bit
				
				bitcount++
			}

			// is the byte full?
			if bitcount == 8 {
				bitcount = 0
				
				if *debug {
					fmt.Printf("out byte: %08b\n", ob)
				}
				
				buf.WriteByte(ob)
				ob = byte(0)
			}
		}
		
		if *debug {
			time.Sleep(1 * time.Second)
		}
	}

	// write leftover
	buf.WriteByte(ob)

	return buf
}

func debiasFile(file string, finfo fs.FileInfo) {

	start := time.Now()
	
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("read", len(data), "bytes from file", file)

	buf := debiasData(data)

	out := filepath.Join(*flagPath, finfo.Name() + "-debiased.bin")
	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	n, err := f.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("wrote", n, "bytes to output file", out, "in", time.Since(start))

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func setBit(n byte, pos uint) byte {
    n |= (1 << pos)
    return n
}