package main

import (
	"fmt"
	"strconv"
)

func clearLine() {
	print("\033[2K\r")
}

// pad the input string up to the given number of space characters
func pad(in string, length int) string {
	if len(in) < length {
		return fmt.Sprintf("%-"+strconv.Itoa(length)+"s", in)
	}
	return in
}