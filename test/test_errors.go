package main

import (
	"errors"

	fdw "github.com/cbandy/go-fdw"
)

type errorsPath string
type errorsIterator struct {
	test string
	row  int
}

func errorsPathError(table string) error {
	if table == "scan_path" {
		return errors.New("bad path")
	}
	return nil
}

func (p errorsPath) Estimate(cost fdw.ScanCostEstimate) (fdw.ScanCostEstimate, error) {
	var err error
	if p == "estimate_scan" {
		err = errors.New("bad estimate")
	}
	return cost, err
}

func (p errorsPath) Begin() (fdw.Iterator, error) {
	var err error
	if p == "begin_scan" {
		err = errors.New("bad begin")
	}
	return &errorsIterator{test: string(p)}, err
}

func (p errorsPath) Close() error {
	if p == "end_path" {
		return errors.New("bad path close")
	}
	return nil
}

func (i errorsIterator) Close() error {
	if i.test == "end_scan" {
		return errors.New("bad iterator close")
	}
	return nil
}

func (i *errorsIterator) Next(attr []fdw.Attribute) (bool, error) {
	i.row++
	if i.test == "during_scan" {
		return i.row < 2, errors.New("bad iteration")
	}
	if i.test == "bad_conversion" {
		return i.row < 2, attr[0].SetString("word")
	}
	return false, nil
}
