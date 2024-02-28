package read

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Execute() float64 {
	result := 0.0
	lines := readChunks("input.txt")
	wss := splitChunk(lines)
	for l := range wss {
		fmt.Printf("%v\n", l)
	}
	return result
}

func readChunks(filename string) chan string {
	out := make(chan string, 100)

	go func() {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		rd := bufio.NewScanner(f)

		for rd.Scan() {
			line := rd.Text()
			out <- line
		}

		close(out)
	}()

	return out
}

type ws struct {
	city    string
	avgTemp float64
}

func splitChunk(in chan string) chan ws {
	out := make(chan ws, 100)

	go func() {

		for l := range in {
			segs := strings.Split(l, ";")

			if len(segs) != 2 {
				panic("wrong")
			}

			city := segs[0]

			avgTemp, err := strconv.ParseFloat(segs[1], 64)

			if err != nil {
				panic(err)
			}

			result := ws{city, avgTemp}

			out <- result
		}

		close(out)

	}()

	return out
}
