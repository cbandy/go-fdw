package fdw

// https://www.postgresql.org/docs/current/static/ddl-foreign-data.html

type Handler interface {
	// called before GetForeignRelSize and GetForeignPaths
	// sufficient before full scan
	EstimateScan(options map[string]string, cost *ScanCostEstimate, private *interface{}) // TODO quals; tlist?

	Scan(
		options map[string]string, // could go away if require user to store in private
		// the plan (scan, thus far)
	) Iterator
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

type ScanCostEstimate struct {
	Rows    float64 // number of tuples
	Width   int     // average width of tuples
	Startup float64
	Total   float64
}
