package fdw

/*
#include "go_fdw.h"
#include "utils/memutils.h"

void
goGetForeignRelSizeWrapper(
	PlannerInfo *root, RelOptInfo *baserel, Oid foreigntableid)
{
	ErrorData *edata = goGetForeignRelSize(root, baserel, foreigntableid);
	if (edata) ReThrowError(edata);
}

void
goBeginForeignScanWrapper(ForeignScanState *node, int eflags)
{
	ErrorData *edata = goBeginForeignScan(node, eflags);
	if (edata) ReThrowError(edata);
}

TupleTableSlot *
goIterateForeignScanWrapper(ForeignScanState *node)
{
	struct goIterateForeignScanResult result = goIterateForeignScan(node);
	if (result.edata) ReThrowError(result.edata);
	return result.slot;
}

void
goEndForeignScanWrapper(ForeignScanState *node)
{
	ErrorData *edata = goEndForeignScan(node);
	if (edata) ReThrowError(edata);
}

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

	fr.GetForeignRelSize = (C.GetForeignRelSize_function)(C.goGetForeignRelSizeWrapper)
	fr.GetForeignPaths = (C.GetForeignPaths_function)(C.goGetForeignPaths)
	fr.GetForeignPlan = (C.GetForeignPlan_function)(C.goGetForeignPlan)
	fr.BeginForeignScan = (C.BeginForeignScan_function)(C.goBeginForeignScanWrapper)
	fr.IterateForeignScan = (C.IterateForeignScan_function)(C.goIterateForeignScanWrapper)
	fr.ReScanForeignScan = (C.ReScanForeignScan_function)(C.goReScanForeignScan)
	fr.EndForeignScan = (C.EndForeignScan_function)(C.goEndForeignScanWrapper)
}
