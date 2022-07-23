package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

type Filter func(int, string) bool

func readLine(r *bufio.Reader) (string, bool, error) {
	s, err := r.ReadString('\n')
	if err == nil {
		return s, true, nil
	}
	if errors.Is(err, io.EOF) {
		if len(s) > 0 {
			return s, true, nil
		}
		return "", false, nil
	}
	return "", false, err
}

func thinout(dst io.Writer, src io.Reader, filter Filter) error {
	w := bufio.NewWriter(dst)
	defer w.Flush()
	r := bufio.NewReader(src)
	readLnum := 0
	for {
		readLnum++
		s, cont, err := readLine(r)
		if err != nil {
			return fmt.Errorf("read failed at line %d: %w", readLnum, err)
		}
		if !cont {
			return nil
		}
		if filter(readLnum, s) {
			_, err := w.WriteString(s)
			if err != nil {
				log.Printf("write failed at input-line %d: %s", readLnum, err)
				return nil
			}
		}
	}
}

func main() {
	var fixes []int
	denominator := flag.Float64("d", 10.0, "denominator of the rate for thin out")
	flag.Func("f", "specify fixed line number (allow multi-times)", func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if n <= 0 {
			return fmt.Errorf("line number to fix must be >= 1, but got %d", n)
		}
		fixes = append(fixes, n)
		return nil
	})
	flag.Parse()

	rate := 1.0 / *denominator
	sort.Ints(fixes)

	err := thinout(os.Stdout, os.Stdin, func(n int, s string) bool {
		if len(fixes) > 0 && n == fixes[0] {
			fixes = fixes[1:]
			return true
		}
		return rand.Float64() < rate
	})
	if err != nil {
		log.Fatal(err)
	}
}
