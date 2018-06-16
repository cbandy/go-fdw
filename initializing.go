package fdw

/*
#include "postgres.h"
#include "foreign/fdwapi.h"

#include "utils/memutils.h"

void goGetForeignRelSize
(PlannerInfo *root, RelOptInfo *baserel, Oid foreigntableid);

void goGetForeignPaths
(PlannerInfo *root, RelOptInfo *baserel, Oid foreigntableid);

ForeignScan *goGetForeignPlan
(PlannerInfo *root, RelOptInfo *baserel, Oid foreigntableid,
 ForeignPath *best_path, List *tlist, List *scan_clauses, Plan *outer_plan);

void goBeginForeignScan
(ForeignScanState *node, int eflags);

TupleTableSlot *goIterateForeignScan
(ForeignScanState *node);

void goReScanForeignScan
(ForeignScanState *node);

void goEndForeignScan
(ForeignScanState *node);

#cgo LDFLAGS: -shared
*/
import "C"
import (
	"log"
	"reflect"
	"unsafe"
)

var initialized struct {
	fdw     FDW
	handler Handler
}

func Initialize(fdw FDW, fdwRoutine interface{}) {
	log.Print("Initialize")

	if C.CurrentMemoryContext != C.MessageContext {
		panic("Unexpected memory context")
	}
	if initialized.handler != nil {
		panic("Handler already assigned")
	}

	initialized.fdw = fdw

	// https://github.com/golang/go/issues/13467
	// frv := reflect.ValueOf(fdwRoutine).Convert(reflect.TypeOf((*C.FdwRoutine)(nil)))
	fr := (*C.FdwRoutine)(unsafe.Pointer(reflect.ValueOf(fdwRoutine).Pointer()))

	fr.GetForeignRelSize = (C.GetForeignRelSize_function)(C.goGetForeignRelSize)
	fr.GetForeignPaths = (C.GetForeignPaths_function)(C.goGetForeignPaths)
	fr.GetForeignPlan = (C.GetForeignPlan_function)(C.goGetForeignPlan)
	fr.BeginForeignScan = (C.BeginForeignScan_function)(C.goBeginForeignScan)
	fr.IterateForeignScan = (C.IterateForeignScan_function)(C.goIterateForeignScan)
	fr.ReScanForeignScan = (C.ReScanForeignScan_function)(C.goReScanForeignScan)
	fr.EndForeignScan = (C.EndForeignScan_function)(C.goEndForeignScan)
}
