package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/DavidEsdrs/1brc/read"
	"github.com/DavidEsdrs/1brc/write"
)

func main() {
	var choose string
	var size int

	flag.IntVar(&size, "size", 1000, "-size 100")
	flag.StringVar(&choose, "choose", "", "read")

	flag.Parse()

	var res float64

	start := time.Now()

	switch {
	case choose == "read":
		res = read.Execute()
	case choose == "write":
		write.Write(size)
	}

	duration := time.Since(start)
	fmt.Printf("%f - %v", res, duration.String())
}
