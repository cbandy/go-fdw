package main

// #include "postgres.h"
// #include "foreign/fdwapi.h"
import "C"
import (
	"errors"

	"github.com/cbandy/go-fdw"
)

func main() {}

//export goTestInitialize
func goTestInitialize(fr *C.FdwRoutine) { fdw.Initialize(entry{}, fr) }

type entry struct{}

func (entry) New() fdw.Handler { return handler{} }

type handler struct{}

func (handler) Scan(t fdw.Table) (fdw.ScanPath, error) {
	var (
		err     error
		name    string
		options map[string]string
	)

	t.Relation(func(r fdw.Relation) { name = r.Name() })
	t.Options(func(o fdw.Options) {
		switch name {
		case "server_options":
			options = o.Server()

		case "table_options":
			options = o.Table()

		case "user_options":
			options, err = o.User()

		default:
			options = o.Server()
			for k, v := range o.Table() {
				options[k] = v
			}
		}
	})

	if err != nil {
		return nil, err
	}

	switch options["test"] {
	case "errors":
		if name == "scan_path" {
			return nil, errors.New("bad path")
		}
		return errorsPath(name), nil
	case "fdw_join":
		return joinsPath(name), nil
	}

	switch name {
	case "server_options", "table_options", "user_options":
		return optionsPath(options), nil
	}

	return nil, errors.New("unexpected test")
}

type optionsPath map[string]string
type optionsIterator struct {
	keys   []string
	values []string
	row    int
}

func (optionsPath) Estimate(cost fdw.ScanCostEstimate) (fdw.ScanCostEstimate, error) { return cost, nil }
func (o optionsPath) Begin() (fdw.Iterator, error) {
	i := optionsIterator{}
	for k, v := range o {
		i.keys = append(i.keys, k)
		i.values = append(i.values, v)
	}
	return &i, nil
}
func (optionsPath) Close() error { return nil }

func (optionsIterator) Close() error { return nil }
func (i *optionsIterator) Next(attr []fdw.Attribute) (bool, error) {
	if i.row >= len(i.keys) {
		return false, nil
	}

	attr[0].SetText([]byte(i.keys[i.row]))
	attr[1].SetText([]byte(i.values[i.row]))
	i.row++
	return true, nil
}

type joinsPath string
type joinsIterator struct {
	values [][]string
	row    int
}

func (joinsPath) Estimate(cost fdw.ScanCostEstimate) (fdw.ScanCostEstimate, error) { return cost, nil }
func (j joinsPath) Begin() (fdw.Iterator, error) {
	switch j {
	case "one":
		return &joinsIterator{values: [][]string{{"1", "x"}, {"2", "y"}, {"3", "z"}}}, nil
	case "two":
		return &joinsIterator{values: [][]string{{"2", "k"}, {"3", "l"}, {"4", "m"}, {"5", "n"}}}, nil
	}
	return nil, nil
}
func (joinsPath) Close() error { return nil }

func (joinsIterator) Close() error { return nil }
func (i *joinsIterator) Next(attr []fdw.Attribute) (bool, error) {
	if i.row >= len(i.values) {
		return false, nil
	}

	for j, a := range attr {
		a.SetString(i.values[i.row][j])
	}
	i.row++
	return true, nil
}
