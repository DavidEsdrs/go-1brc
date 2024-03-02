package main

import (
	"fmt"
	"time"

	"github.com/DavidEsdrs/1brc/read"
)

func main() {
	var res float64

	start := time.Now()

	res = read.Execute()

	duration := time.Since(start)
	fmt.Printf("%f - %v", res, duration.String())
}
