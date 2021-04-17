package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
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
	flagKaminsky = flag.Bool("k", false, "kaminsky debiasing")
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

// The Von Neumann Debiasing algorithm works on pairs of bits, and produces output as follows:
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

// The Von Neumann Debiasing algorithm works on pairs of bits, and produces output as follows:
// - If the input is "00" or "11", the input is discarded (no output).
// - If the input is "10", output a "1".
// - If the input is "01", output a "0".
// 
// Kaminsky:
// - collect discarded bytes
// - use discarded bytes as input for SHA512
// - use the SHA512 hash as key for encrypting the output data with AES
func debiasDataKaminsky(data []byte) (buf bytes.Buffer, discardBuf bytes.Buffer) {
	var (
		
		// discard byte
		db = byte(0)
		discardBitcount uint

		// out byte
		ob = byte(0)
		bitcount uint
	)

	for _, b := range data {

		if *debug {
			fmt.Printf("byte: %08b\n", b)
		}
	
		for j:=0; j<8; j+= 2 {
			
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
			} else {
				// discarded bits: collect
				db = setBit(db, 7-discardBitcount)
				discardBitcount++
			}

			if discardBitcount == 8 {
				discardBitcount = 0

				if *debug {
					fmt.Printf("discard out byte: %08b\n", db)
				}
				
				discardBuf.WriteByte(db)
				db = byte(0)
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
	discardBuf.WriteByte(db)

	return
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func debiasFile(file string, finfo fs.FileInfo) {

	start := time.Now()
	
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("read", len(data), "bytes from file", file)

	var buf bytes.Buffer
	if *flagKaminsky {
		var discardBuf bytes.Buffer
		buf, discardBuf = debiasDataKaminsky(data)

		fmt.Println("kaminsky mode, discard buffer:", discardBuf.Len())

		// create SHA256
		h := sha256.Sum256(discardBuf.Bytes())
		
		// convert into []byte to please go compiler
		var key = make([]byte, 32)
		for i:=0; i<32; i++ {
			key[i] = h[i]
		}

		// TODO: read value from dev random. for our experiments, it better to have a static value here.
		iv := []byte("1234567890123456")

		// pad plaintext
		bPlaintext := PKCS5Padding(buf.Bytes(), aes.BlockSize, buf.Len())

		// init cipher with key
		block, err := aes.NewCipher(key)
		if err != nil {
			log.Fatal(err)
		}
		
		ciphertext := make([]byte, len(bPlaintext))
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(ciphertext, bPlaintext)

		buf.Reset()
		buf.Write(ciphertext)
	} else {
		buf = debiasData(data)
	}
	
	var out string
	if *flagKaminsky {
		out = filepath.Join(*flagPath, finfo.Name() + "-ka-debiased.bin")
	} else {
		out = filepath.Join(*flagPath, finfo.Name() + "-vn-debiased.bin")
	}
	
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