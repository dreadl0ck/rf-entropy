package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	debug = false
	path = "samples"
)

func main() {
	
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		
		if filepath.Ext(f.Name()) != ".wav" {
			continue
		}
		
		file := filepath.Join(path, f.Name())
		fmt.Println("processing", file)
		
		debiasFile(file, f)
	}
}

func debiasData(data []byte) bytes.Buffer {
	var (
		buf bytes.Buffer
		ob = byte(0)
		bitcount uint
	)

	for _, b := range data {

		if debug {
			fmt.Printf("byte: %08b\n", b)
		}
	
		for j:=0; j<8; j+= 2 {

			ch := (b >> j) & 0x01
			ch2 := (b >> (j+1)) & 0x01

			if debug {
				fmt.Println(ch, ch2)
			}
			
			if (ch != ch2) {
			
				if ch != 1 {
					/* store a 1 in our bitbuffer */
					ob = setBit(ob, bitcount)

					if debug {
						fmt.Printf("collecting 1, out byte: %08b\n", ob)
					}
				} /* else, leave the buffer alone, it's already 0 at this bit */
				
				bitcount++
			}
		} 

		/* is byte full? */
		if bitcount == 8 {
			bitcount = 0
			
			if debug {
				fmt.Printf("out byte: %08b\n", ob)
			}
			
			buf.WriteByte(ob)
			ob = byte(0)
		}
		
		if debug {
			time.Sleep(1 * time.Second)
		}
	}

	return buf
}

func debiasFile(file string, finfo fs.FileInfo) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	buf := debiasData(data)

	out := filepath.Join(path, finfo.Name() + "-debiased.bin")
	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	n, err := f.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("wrote", n, "bytes to output file", out)

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func setBit(n byte, pos uint) byte {
    n |= (1 << pos)
    return n
}