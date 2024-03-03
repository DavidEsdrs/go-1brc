package read

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

func Execute() float64 {
	result := 0.0
	lines := readChunks("input.txt")
	fixedChunks := fixChunks(lines)
	segs := splitChunk(fixedChunks)
	wss := parseSegments(segs)
	stations := processStation(wss)

	stations = mergeSort(stations)

	for _, s := range stations {
		fmt.Printf("%v=%v/%v/%v\n", s.city, s.minT, s.avgT, s.maxT)
	}

	return result
}

func readChunks(filename string) chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		buffer := make([]byte, 1024*4)

		rd := bufio.NewReader(f)

		for {
			n, err := rd.Read(buffer)

			if err != nil && err == io.EOF {
				break
			}

			out <- string(buffer[:n])
		}
	}()

	return out
}

func fixChunks(in chan string) chan string {
	out := make(chan string, 1000)

	go func() {
		defer close(out)

		var lastOverflow string
		var builder strings.Builder

		for chunk := range in {
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
			builder.Reset()
		}
	}()

	return out
}

type ws struct {
	city             string
	minT, maxT, avgT float64
	count            int
}

func splitChunk(in chan string) chan [2]string {
	out := make(chan [2]string, 1000)

	go func() {
		defer close(out)

		var result [2]string = [2]string{}

		for l := range in {
			segs := strings.Split(l, ";")

			if len(segs) != 2 {
				panic("wrong seg count")
			}

			result[0] = segs[0]
			result[1] = segs[1]

			out <- result
		}
	}()

	return out
}

func parseSegments(in chan [2]string) chan *ws {
	out := make(chan *ws, 1000)

	go func() {
		defer close(out)

		for segs := range in {
			avgTemp, err := strconv.ParseFloat(segs[1], 64)

			if err != nil {
				panic(err)
			}

			result := ws{city: segs[0], avgT: avgTemp}

			out <- &result
		}
	}()

	return out
}

func processStation(in chan *ws) []ws {
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

	var list []ws

	for _, s := range result {
		list = append(list, s)
	}

	return list
}

func mergeSort(items []ws) []ws {
	if len(items) < 2 {
		return items
	}

	mid := len(items) / 2
	left := mergeSort(items[:mid])
	right := mergeSort(items[mid:])

	return merge(left, right)
}

func merge(a, b []ws) []ws {
	result := make([]ws, 0)
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i].city[0] < b[j].city[0] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}

	// Append any remaining elements from both slices
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)

	return result
}
