package main

import "C"
import (
	"fmt"
	"runtime"

	fdw "github.com/cbandy/go-fdw"
)

var benchmark struct {
	mem runtime.MemStats
}

//export goTestBenchmarkBegin
func goTestBenchmarkBegin() {
	runtime.ReadMemStats(&benchmark.mem)
}

//export goTestBenchmarkEnd
func goTestBenchmarkEnd(dest []byte) C.int {
	var end runtime.MemStats
	runtime.ReadMemStats(&end)

	s := fmt.Sprintf("%8d B\t%8d allocs",
		end.TotalAlloc-benchmark.mem.TotalAlloc,
		end.Mallocs-benchmark.mem.Mallocs)

	return C.int(copy(dest, s))
}

type benchmarkPath string
type benchmarkIterator struct {
	data [][]string
	row  int
}

var benchmarkConst = [][]string{
	{"one", "a", "b", "Lorem ipsum dolor sit amet"},
	{"two", "c", "d", "consectetur adipiscing elit"},
	{"three", "e", "f", "sed do eiusmod tempor incididunt"},
	{"four", "g", "h", "ut labore et dolore magna aliqua"},
}

func (p benchmarkPath) Estimate(cost fdw.ScanCostEstimate) (fdw.ScanCostEstimate, error) {
	return cost, nil
}

func (p benchmarkPath) Begin() (fdw.Iterator, error) {
	return &benchmarkIterator{data: benchmarkConst}, nil
}

func (p benchmarkPath) Close() error { return nil }

func (i *benchmarkIterator) Next(attr []fdw.Attribute) (bool, error) {
	if i.row < len(i.data) {
		for j, a := range attr {
			a.SetString(i.data[i.row][j])
		}
		i.row++
		return true, nil
	}
	return false, nil
}

func (benchmarkIterator) Close() error { return nil }
