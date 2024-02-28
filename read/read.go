package read

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

func Execute() float64 {
	result := 0.0
	lines := readChunks("input.txt")
	wss := splitChunk(lines)
	// stations := processStation(wss)
	processStation(wss)
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
	city             string
	minT, maxT, avgT float64
	count            int
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

			result := ws{city: city, avgT: avgTemp}

			out <- result
		}

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
