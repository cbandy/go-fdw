package fdw

// https://www.postgresql.org/docs/current/static/ddl-foreign-data.html

type Handler interface {
	Scan(table Table) ScanPath
}

//type CostEstimate struct {
//	Rows    float64 // number of tuples
//	Startup float64
//	Total   float64
//}

type Attribute interface {
	// FIXME terrible names
	SetText([]byte)  // error?
	SetText0([]byte) // error?
	SetText2(string) // error?
	TypeOid() uint
	// type name?
}

type Iterator interface {
	Close() // TODO error
	HasNext() bool
	Next([]Attribute) // TODO error
	// TODO rescan (parameters could change)
}

type Options interface {
	Server() map[string]string
	Table() map[string]string
	User() map[string]string
}

type Relation interface {
	Name() string
}

type ScanCostEstimate struct {
	Rows    float64 // number of tuples
	Width   int     // average width of tuples
	Startup float64
	Total   float64
}

type ScanPath interface {
	Estimate(ScanCostEstimate) ScanCostEstimate
	Begin() Iterator
	Close()
}

type Table interface {
	Oid() uint
	Options(func(Options))
	Relation(func(Relation))
}
