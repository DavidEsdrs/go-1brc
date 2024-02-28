package read

import (
	"bufio"
	"os"
)

func Execute() float64 {
	result := 0.0
	lines := readChunks("input.txt")
	for l := range lines {
		println(l)
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
