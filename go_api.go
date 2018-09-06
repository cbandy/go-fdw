package fdw

// https://www.postgresql.org/docs/current/static/ddl-foreign-data.html

type FDW interface {
	New() Handler
}

type Handler interface {
	Scan(table Table) (ScanPath, error)
}

//type CostEstimate struct {
//	Rows    float64 // number of tuples
//	Startup float64
//	Total   float64
//}

type Attribute interface {
	SetText([]byte) error   // same as SetText0, but slice is copied first to NUL-terminate
	SetText0([]byte) error  //
	SetString(string) error // same as SetText0, but string is copied first to NUL-terminate
	TypeOid() uint
	// type name?
}

type Iterator interface {
	Close() error
	Next([]Attribute) (bool, error)
	// TODO rescan (parameters could change)
}

type Options interface {
	Server() map[string]string
	Table() map[string]string
	User() (map[string]string, error) // returns an error when the user mapping does not exist
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
	Estimate(ScanCostEstimate) (ScanCostEstimate, error)
	Begin() (Iterator, error)
	Close() error
}

type Table interface {
	Oid() uint
	Options(func(Options))
	Relation(func(Relation))
}
