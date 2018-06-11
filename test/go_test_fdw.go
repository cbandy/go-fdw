package main

// #include "postgres.h"
// #include "foreign/fdwapi.h"
import "C"
import "github.com/cbandy/go-fdw"

func main() {}

//export goTestInitialize
func goTestInitialize(fr *C.FdwRoutine) { fdw.Initialize(handler{}, fr) }

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

	//switch name {
	//case "server_options", "table_options", "user_options":
	return optionsPath(options)
	//}
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
