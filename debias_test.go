package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
)

type bitString string

func (b bitString) AsByteSlice() []byte {
    var out []byte
    var str string

    for i := len(b); i > 0; i -= 8 {
        if i-8 < 0 {
            str = string(b[0:i])
        } else {
            str = string(b[i-8 : i])
        }
        v, err := strconv.ParseUint(str, 2, 8)
        if err != nil {
            panic(err)
        }
        out = append([]byte{byte(v)}, out...)
    }
    return out
}

func (b bitString) AsHexSlice() []string {
    var out []string
    byteSlice := b.AsByteSlice()
    for _, b := range byteSlice {
        out = append(out, "0x" + hex.EncodeToString([]byte{b}))
    }
    return out
}

func TestBitStringConversion(t *testing.T) {
	if !bytes.Equal(bitString("00000000").AsByteSlice(), []byte{byte(0)}) {
		t.Fatal("incorrect conversion result: ", bitString("00000000").AsByteSlice(), " expected: ", []byte{byte(0)})
	}
	if !bytes.Equal(bitString("11111111").AsByteSlice(), []byte{byte(255)}) {
		t.Fatal("incorrect conversion result: ", bitString("11111111").AsByteSlice(), " expected: ", []byte{byte(255)})
	}
	if !bytes.Equal(bitString("10101010").AsByteSlice(), []byte{byte(170)}) {
		t.Fatal("incorrect conversion result: ", bitString("10101010").AsByteSlice(), " expected: ", []byte{byte(170)})
	}
	if !bytes.Equal(bitString("01010101").AsByteSlice(), []byte{byte(85)}) {
		t.Fatal("incorrect conversion result: ", bitString("01010101").AsByteSlice(), " expected: ", []byte{byte(85)})
	}
}

func TestVonNeumannDebiasing(t *testing.T) {

	// 0 to 0000000000000000000 (steps of one) --> nothing
	buf := debiasData(bitString("0000000000000000000").AsByteSlice())
	if len(buf.Bytes()) != 0 {
		fmt.Printf("%08b", buf.Bytes())
		t.Fatal("expected no output")
	}

	// 1 to 1111111111111111111 (steps of one) --> nothing
	buf = debiasData(bitString("1111111111111111111").AsByteSlice())
	if len(buf.Bytes()) != 0 {
		fmt.Printf("%08b", buf.Bytes())
		t.Fatal("expected no output")
	}

	debug = true

	// 01 to 01010101 01010101 (steps of two) --> all zeros * 1/2 input length
	buf = debiasData(bitString("0101010101010101").AsByteSlice())
	if len(buf.Bytes()) != 1 || !bytes.Equal(buf.Bytes(), []byte{byte(0)}) {
		fmt.Printf("%08b", buf.Bytes())
		t.Fatal("expected no output")
	}

	// 10 to 10101010 10101010 (steps of two) --> all ones * 1/2 input length
	buf = debiasData(bitString("1010101010101010").AsByteSlice())
	if len(buf.Bytes()) != 1 || !bytes.Equal(buf.Bytes(), []byte{byte(255)}) {
		fmt.Printf("%08b", buf.Bytes())
		t.Fatal("expected no output")
	}
}