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

// Scan dispatches to various test cases
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
		return errorsPath(name), errorsPathError(name)
	case "fdw_join":
		return joinsPath(name), nil
	}

	switch name {
	case "server_options", "table_options", "user_options":
		return optionsPath(options), nil
	}

	return nil, errors.New("unexpected test")
}
