package main

// #include "postgres.h"
// #include "foreign/fdwapi.h"
import "C"
import "github.com/cbandy/go-fdw"

func main() {}

//export goTestInitialize
func goTestInitialize(fr *C.FdwRoutine) { fdw.Initialize(handler{}, fr) }

type handler struct{}

func (handler) EstimateScan(options map[string]string, cost *fdw.ScanCostEstimate, private *interface{}) {
	cost.Rows = 7
	cost.Startup = 100
	cost.Total = 999

	*private = options
}

func (handler) Scan(options map[string]string) fdw.Iterator {
	scanner := optionsScanner{}
	for k, v := range options {
		scanner.keys = append(scanner.keys, k)
		scanner.values = append(scanner.values, v)
	}
	return &scanner
}

type optionsScanner struct {
	keys   []string
	values []string
	row    int
}

func (optionsScanner) Close()          {}
func (s optionsScanner) HasNext() bool { return s.row < len(s.keys) }
func (s *optionsScanner) Next(attr []fdw.Attribute) {
	attr[0].SetText([]byte(s.keys[s.row]))
	attr[1].SetText([]byte(s.values[s.row]))
	s.row++
}
