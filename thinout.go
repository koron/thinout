package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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

type seedOpt struct {
	time bool
	seed int64
}

func (so *seedOpt) String() string {
	if so.time {
		return "time"
	}
	return strconv.FormatInt(so.seed, 10)
}

func (so *seedOpt) Set(s string) error {
	switch s {
	case "time":
		*so = seedOpt{true, 0}
		return nil
	default:
		seed, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*so = seedOpt{false, seed}
		return nil
	}
}

func (so *seedOpt) rand() *rand.Rand {
	seed := so.seed
	if so.time {
		seed = time.Now().UnixMilli()
	}
	return rand.New(rand.NewSource(seed))
}

type intArray []int

func (ia intArray) String() string {
	bb := &bytes.Buffer{}
	for i, v := range ia {
		if i != 0 {
			bb.WriteRune(',')
		}
		bb.WriteString(strconv.Itoa(v))
	}
	return bb.String()
}

func (ia *intArray) Set(s string) error {
	max := strings.Count(s, ",") + 1
	vals := make([]int, 0, max)
	for _, t := range strings.SplitN(s, ",", max) {
		n, err := strconv.Atoi(t)
		if err != nil {
			return fmt.Errorf("number required: %w", err)
		}
		if n <= 0 {
			return fmt.Errorf("line number to fix must be >= 1, but got %d", n)
		}
		vals = append(vals, n)
	}
	*ia = append(*ia, vals...)
	return nil
}

func main() {
	var (
		rate  float64
		seed  seedOpt
		fixes intArray
	)

	flag.Float64Var(&rate, "r", 0.1, "rate to output [0.0,1.0]")
	flag.Var(&seed, "s", `seed for random numbers: "time" or int64 (default: 0)`)
	flag.Var(&fixes, "f", `specify fixed line numbers, comma separated`)
	flag.Parse()

	if rate < 0 || rate > 1 {
		log.Fatalf("output rate (-r) out of range [0.0,1.0]: %e", rate)
	}
	rnd := seed.rand()
	sort.Ints(fixes)

	err := thinout(os.Stdout, os.Stdin, func(n int, s string) bool {
		if len(fixes) > 0 && n == fixes[0] {
			fixes = fixes[1:]
			return true
		}
		return rnd.Float64() < rate
	})
	if err != nil {
		log.Fatal(err)
	}
}
