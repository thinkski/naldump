// Copyright 2023 Chris Hiszpanski
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the “Software”),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	flagExclude string
	flagHevc    bool
)

func init() {
	flag.StringVar(&flagExclude, "exclude", "", " NAL types to exclude (comma separated)")
	flag.BoolVar(&flagHevc, "hevc", false, "Consider bytestream to be HEVC")
}

func main() {
	flag.Parse()

	exclude := make(map[int]bool)
	for _, ign := range strings.Split(flagExclude, ",") {
		if i, err := strconv.Atoi(ign); nil == err {
			exclude[i] = true
		}
	}

	var err error
	var f *os.File

	// open input
	switch flag.NArg() {
	case 0:
		f = os.Stdin
	case 1:
		if f, err = os.Open(flag.Arg(0)); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open file: %s", err.Error())
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unexpected number of arguments")
		os.Exit(1)
	}
	defer f.Close()

	buffer := make([]byte, 256*1024)
	scanner := bufio.NewScanner(f)
	scanner.Buffer(buffer, 8192*1024)
	scanner.Split(splitNalUnit)
	n := 0

	start := time.Now()

	var names map[int]string
	if flagHevc {
		names = map[int]string{
			1:  "TRAIL_R",
			19: "IDR_W_RADL",
			32: "video parameter set",
			33: "sequence parameter set",
			34: "picture parameter set",
		}
	} else {
		names = map[int]string{
			1: "non-idr coded picture",
			5: "    idr coded picture",
			6: "supplemental enhancement information",
			7: "sequence parameter set",
			8: "picture parameter set",
		}
	}

	for scanner.Scan() {
		b := scanner.Bytes()

		var typ int
		if flagHevc {
			typ = int((b[0] >> 1) & 0x3f)
		} else {
			typ = int(b[0] & 0x1f)
		}

		if _, ok := exclude[int(typ)]; ok {
			n++
			continue
		}

		if name, ok := names[typ]; ok {
			if f == os.Stdin {
				elapsed := time.Now().Sub(start).Seconds()
				fmt.Printf("%.6f\t%v\t%v\t%v\t%s\n", elapsed, n, typ, len(b), name)
			} else {
				fmt.Printf("%v\t%v\t%v\t%s\n", n, typ, len(b), name)
			}
		}
		n++
	}
}

// splitNalUnit adheres to bufio.SplitFunc signature and returns next NAL unit
func splitNalUnit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	startPrefix := []byte{0, 0, 1}

	// find first prefix code
	left := bytes.Index(data, startPrefix)

	switch left {
	case 0:
		// no data other than prefix code? need more data.
		if 3 == len(data) {
			return 0, nil, nil
		}
		break
	case -1:
		// prefix code not found. need more data.
		return len(data), nil, nil
	default:
		// advance prefix code to be first byte in scanner buffer
		return left, nil, nil
	}

	// find next prefix code
	right := bytes.Index(data[len(startPrefix):], startPrefix)

	switch right {
	case -1:
		// no second prefix code. need more data.
		return 0, nil, nil
	default:
		// second prefix found. return data between first and second prefix code.
		return len(startPrefix) + right,
			data[len(startPrefix)+left : len(startPrefix)+right-1],
			nil
	}
}
