package main

import fdw "github.com/cbandy/go-fdw"

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

func (i *optionsIterator) Next(attr []fdw.Attribute) (ok bool, err error) {
	if ok = i.row < len(i.keys); ok {
		if err = attr[0].SetText([]byte(i.keys[i.row])); err == nil {
			err = attr[1].SetText([]byte(i.values[i.row]))
		}
		i.row++
	}
	return
}
