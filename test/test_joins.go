package main

import fdw "github.com/cbandy/go-fdw"

type joinsPath string
type joinsIterator struct {
	values [][]string
	row    int
}

func (joinsPath) Estimate(cost fdw.ScanCostEstimate) (fdw.ScanCostEstimate, error) { return cost, nil }

func (p joinsPath) Begin() (fdw.Iterator, error) {
	switch p {
	case "one":
		return &joinsIterator{values: [][]string{{"1", "x"}, {"2", "y"}, {"3", "z"}}}, nil
	case "two":
		return &joinsIterator{values: [][]string{{"2", "k"}, {"3", "l"}, {"4", "m"}, {"5", "n"}}}, nil
	}
	return nil, nil
}

func (joinsPath) Close() error { return nil }

func (joinsIterator) Close() error { return nil }

func (i *joinsIterator) Next(attr []fdw.Attribute) (ok bool, err error) {
	if ok = i.row < len(i.values); ok {
		for j, a := range attr {
			if err = a.SetString(i.values[i.row][j]); err != nil {
				break
			}
		}
		i.row++
	}
	return
}
