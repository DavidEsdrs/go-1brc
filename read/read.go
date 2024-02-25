package read

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

func Execute() float64 {
	result := 0.0
	chunks := readChunk()
	fixedChunks := capUntilLastEndOfLine(chunks)
	lines := parseLines(fixedChunks)
	floats := processLine(lines)

	for f := range floats {
		result += f
	}
	return result
}

func readChunk() chan []byte {
	out := make(chan []byte, 100)
	go func() {
		f, err := os.Open("input.txt")

		if err != nil {
			panic(err)
		}

		defer f.Close()

		b := bufio.NewReader(f)

		buffer := make([]byte, 1024*256)

		for {
			_, err := b.Read(buffer)

			if err != nil && errors.Is(err, io.EOF) {
				break
			}

			out <- buffer

		}

		close(out)
	}()
	return out
}

func capUntilLastEndOfLine(in chan []byte) chan []byte {
	out := make(chan []byte, 100)
	go func() {
		for chunk := range in {
			length := len(chunk)
			for i := length - 1; i > 0; i-- {
				if chunk[i] == '\n' {
					chunk = chunk[:i]
					break
				}
			}
			out <- chunk
		}
		close(out)
	}()
	return out
}

func parseLines(in chan []byte) chan string {
	out := make(chan string, 100)

	var wg sync.WaitGroup

	process := func(chunk []byte) {
		defer wg.Done()
		var builder strings.Builder

		for _, b := range chunk {
			if b == '\n' {
				str := builder.String()
				builder.Reset()
				out <- str
			} else {
				builder.WriteByte(b)
			}
		}
	}

	go func() {
		for chunk := range in {
			wg.Add(1)
			go process(chunk)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func processLine(in chan string) chan float64 {
	out := make(chan float64)

	var wg sync.WaitGroup

	process := func(line string) {
		defer wg.Done()

		for i, b := range line {
			if b == ';' {
				floatStr := line[i+1:]
				f, err := strconv.ParseFloat(floatStr, 64)
				if err != nil {
					break
				}
				out <- f
				break
			}
		}
	}

	go func() {
		for line := range in {
			wg.Add(1)
			go process(line)
		}

		wg.Wait()
		close(out)
	}()

	return out
}
