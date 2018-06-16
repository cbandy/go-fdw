package main

// #include "postgres.h"
// #include "foreign/fdwapi.h"
import "C"
import "github.com/cbandy/go-fdw"

func main() {}

//export goTestInitialize
func goTestInitialize(fr *C.FdwRoutine) { fdw.Initialize(entry{}, fr) }

type entry struct{}

func (entry) New() fdw.Handler { return handler{} }

type handler struct{}

func (handler) Scan(t fdw.Table) fdw.ScanPath {
	var name string
	t.Relation(func(r fdw.Relation) { name = r.Name() })

	var options map[string]string
	t.Options(func(o fdw.Options) {
		switch name {
		case "server_options":
			options = o.Server()

		case "table_options":
			options = o.Table()

		case "user_options":
			options = o.User()

		default:
			options = o.Server()
			for k, v := range o.Table() {
				options[k] = v
			}
		}
	})

	switch options["test"] {
	case "fdw_join":
		return joinsPath(name)
	}

	switch name {
	case "server_options", "table_options", "user_options":
		return optionsPath(options)
	}

	return nil
}

type optionsPath map[string]string
type optionsIterator struct {
	keys   []string
	values []string
	row    int
}

func (optionsPath) Estimate(cost fdw.ScanCostEstimate) fdw.ScanCostEstimate { return cost }
func (o optionsPath) Begin() fdw.Iterator {
	i := optionsIterator{}
	for k, v := range o {
		i.keys = append(i.keys, k)
		i.values = append(i.values, v)
	}
	return &i
}
func (optionsPath) Close() {}

func (optionsIterator) Close()          {}
func (i optionsIterator) HasNext() bool { return i.row < len(i.keys) }
func (i *optionsIterator) Next(attr []fdw.Attribute) {
	attr[0].SetText([]byte(i.keys[i.row]))
	attr[1].SetText([]byte(i.values[i.row]))
	i.row++
}

type joinsPath string
type joinsIterator struct {
	values [][]string
	row    int
}

func (joinsPath) Estimate(cost fdw.ScanCostEstimate) fdw.ScanCostEstimate { return cost }
func (j joinsPath) Begin() fdw.Iterator {
	switch j {
	case "one":
		return &joinsIterator{values: [][]string{{"1", "x"}, {"2", "y"}, {"3", "z"}}}
	case "two":
		return &joinsIterator{values: [][]string{{"2", "k"}, {"3", "l"}, {"4", "m"}, {"5", "n"}}}
	}
	return nil
}
func (joinsPath) Close() {}

func (joinsIterator) Close()          {}
func (i joinsIterator) HasNext() bool { return i.row < len(i.values) }
func (i *joinsIterator) Next(attr []fdw.Attribute) {
	for j, a := range attr {
		a.SetString(i.values[i.row][j])
	}
	i.row++
}
