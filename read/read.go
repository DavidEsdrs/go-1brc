package read

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

func Execute() float64 {
	result := 0.0
	lines := readChunks("input.txt")
	fixedChunks := fixChunks(lines)
	wss := splitChunk(fixedChunks)
	stations := processStation(wss)
	processStation(wss)
	for _, v := range stations {
		fmt.Printf("%#v\n", v)
	}
	return result
}

func readChunks(filename string) chan string {
	out := make(chan string)

	go func() {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		bufferSize := 57
		// buffer := make([]byte, 1024*256)
		buffer := make([]byte, bufferSize)

		rd := bufio.NewReader(f)

		for {
			n, err := rd.Read(buffer)

			if err != nil && err == io.EOF {
				break
			}

			out <- string(buffer[:n])
		}

		close(out)
	}()

	return out
}

func fixChunks(in chan string) chan string {
	out := make(chan string, 100)

	go func() {
		defer close(out)

		var lastOverflow string

		for chunk := range in {
			var builder strings.Builder

			builder.WriteString(lastOverflow)

			for i := 0; i < len(chunk); i++ {
				if chunk[i] == '\n' {
					str := builder.String()
					builder.Reset()
					if str != "" {
						out <- str
					}
				} else {
					builder.WriteByte(chunk[i])
				}
			}

			lastOverflow = builder.String()
		}
	}()

	return out
}

type ws struct {
	city             string
	minT, maxT, avgT float64
	count            int
}

func splitChunk(in chan string) chan ws {
	out := make(chan ws, 100)

	var wg sync.WaitGroup

	go func() {

		for l := range in {
			wg.Add(1)

			go func(l string) {
				defer wg.Done()

				segs := strings.Split(l, ";")

				if len(segs) != 2 {
					panic("wrong seg count")
				}

				city := segs[0]

				avgTemp, err := strconv.ParseFloat(segs[1], 64)

				if err != nil {
					panic(err)
				}

				result := ws{city: city, avgT: avgTemp}

				out <- result
			}(l)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func processStation(in chan ws) map[string]ws {
	result := make(map[string]ws)

	for s := range in {
		if r, ok := result[s.city]; ok {
			r.avgT += s.avgT

			if s.avgT < r.minT {
				r.minT = s.avgT
			}

			if s.avgT > r.maxT {
				r.maxT = s.avgT
			}

			r.count++

			result[s.city] = r
		} else {
			result[s.city] = ws{
				city:  s.city,
				avgT:  s.avgT,
				minT:  s.avgT,
				maxT:  s.avgT,
				count: 1,
			}
		}
	}

	for k, v := range result {
		v.avgT /= float64(v.count)
		v.avgT = math.Ceil(v.avgT)
		result[k] = v
	}

	return result
}
