package fdw

/*
#include "postgres.h"
#include "foreign/fdwapi.h"

#include "funcapi.h"
#include "nodes/makefuncs.h"
#include "optimizer/pathnode.h"


#include "postgres.h"
#include "access/htup_details.h"
#include "access/reloptions.h"
#include "access/sysattr.h"
#include "catalog/pg_foreign_table.h"
#include "commands/copy.h"
#include "commands/defrem.h"
#include "commands/explain.h"
#include "commands/vacuum.h"
#include "foreign/fdwapi.h"
#include "foreign/foreign.h"
#include "funcapi.h"
#include "miscadmin.h"
#include "nodes/makefuncs.h"
#include "optimizer/cost.h"
#include "optimizer/pathnode.h"
#include "optimizer/planmain.h"
#include "optimizer/restrictinfo.h"
#include "optimizer/var.h"
#include "utils/memutils.h"
#include "utils/rel.h"
#include "utils/palloc.h"

void goGetForeignRelSize
(PlannerInfo *root, RelOptInfo *baserel, Oid foreigntableid);

void goGetForeignPaths
(PlannerInfo *root, RelOptInfo *baserel, Oid foreigntableid);

ForeignScan *goGetForeignPlan
(PlannerInfo *root, RelOptInfo *baserel, Oid foreigntableid,
 ForeignPath *best_path, List *tlist, List *scan_clauses, Plan *outer_plan);

void goBeginForeignScan              (ForeignScanState *node, int eflags);
TupleTableSlot *goIterateForeignScan (ForeignScanState *node);
void goReScanForeignScan             (ForeignScanState *node);
void goEndForeignScan                (ForeignScanState *node);

static inline void
goAddForeignScanPath(
	PlannerInfo *root, RelOptInfo *baserel,
	PathTarget *target,
	double rows, Cost startup_cost, Cost total_cost,
	List *pathkeys,
	Relids required_outer,
	Path *fdw_outerpath,
	List *fdw_private)
{
	add_path(baserel, (Path *)create_foreignscan_path(
		root, baserel, target,
		rows, startup_cost, total_cost,
		pathkeys, required_outer, fdw_outerpath,
		fdw_private));
}

static inline void
goClearTupleSlot(TupleTableSlot *slot)
{
	ExecClearTuple(slot);
	memset(slot->tts_values, 0, sizeof(Datum) * slot->tts_tupleDescriptor->natts);
	memset(slot->tts_isnull, true, sizeof(bool) * slot->tts_tupleDescriptor->natts);
}

static inline List *
goExtractBareClauses(List *restrictinfo_list)
{
	return extract_actual_clauses(restrictinfo_list, false);
}

static inline Oid
goSlotGetTypeOid(TupleTableSlot *slot, int i)
{
	return slot->tts_tupleDescriptor->attrs[i]->atttypid;
}

static inline void
goSlotSetText(TupleTableSlot *slot, AttInMetadata *attinmeta, int i, char *value)
{
	slot->tts_isnull[i] = (value == NULL) ? true : false;
	slot->tts_values[i] = InputFunctionCall(
		&attinmeta->attinfuncs[i], value, attinmeta->attioparams[i], attinmeta->atttypmods[i]);
}

#cgo LDFLAGS: -shared
*/
import "C"
import (
	"log"
	"reflect"
	"unsafe"
)

type attribute struct {
	index C.int
	state *execState
}

func (a attribute) SetText(input []byte) {
	buffer := make([]byte, len(input)+1)
	buffer[len(input)] = 0
	copy(buffer, input)
	a.SetText0(buffer)
}

func (a attribute) SetText0(input []byte) {
	C.goSlotSetText(a.state.slot, a.state.attinmeta, a.index, (*C.char)(unsafe.Pointer(&input[0])))
}

func (a attribute) SetText2(input string) {
	buffer := make([]byte, len(input)+1)
	buffer[len(input)] = 0
	copy(buffer, input)
	a.SetText0(buffer)
}

func (a attribute) TypeOid() uint {
	return uint(C.goSlotGetTypeOid(a.state.slot, a.index))
}

type execState struct {
	scanCost ScanCostEstimate
	scanPath ScanPath

	iterator Iterator

	attributes []Attribute
	attinmeta  *C.AttInMetadata
	slot       *C.TupleTableSlot
}

func makeOptions(defElems *C.List) map[string]string {
	result := make(map[string]string)

	if defElems != nil {
		for lc := defElems.head; lc != nil; lc = lc.next {
			de := *(**C.DefElem)(unsafe.Pointer(&lc.data[0]))
			result[C.GoString(de.defname)] = C.GoString(C.defGetString(de))
		}
	}

	log.Printf("%v", result)

	return result
}

func Initialize(handler Handler, fdwRoutine interface{}) {
	initialized.handler = handler

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

// https://www.postgresql.org/docs/current/static/fdw-callbacks.html#FDW-CALLBACKS-SCAN
//
// Planning

//export goGetForeignRelSize
func goGetForeignRelSize(root *C.PlannerInfo, baserel *C.RelOptInfo, foreigntableid C.Oid) {
	log.Printf("%p GetForeignRelSize", root)

	state := new(execState)
	state.scanPath = initialized.handler.Scan(table{foreigntableid})
	state.scanCost = state.scanPath.Estimate(ScanCostEstimate{
		Rows:  float64(baserel.rows),
		Width: int(baserel.reltarget.width),
	})

	baserel.rows = C.double(state.scanCost.Rows)
	baserel.reltarget.width = C.int(state.scanCost.Width)

	initialized.execs.Store(root, state)
}

//export goGetForeignPaths
func goGetForeignPaths(root *C.PlannerInfo, baserel *C.RelOptInfo, foreigntableid C.Oid) {
	log.Printf("%p GetForeignPaths", root)

	state := initialized.execs.Load(root)

	// 1. full scan
	C.goAddForeignScanPath(root,
		baserel, nil,
		baserel.rows,
		C.Cost(state.scanCost.Startup),
		C.Cost(state.scanCost.Total),
		nil, nil, nil, nil /* private */)

	// 2. sorted result
	// ask the handler if it is okay with _all_ the pathkeys
	// if so, add another path with (same, no-copy) pathkeys populated
	// same row count, new startup and total costs

	// 3. parameterized results
	// add a path (rows and costs) for each join condition
}

//export goGetForeignPlan
func goGetForeignPlan(root *C.PlannerInfo, baserel *C.RelOptInfo, foreigntableid C.Oid,
	best_path *C.ForeignPath, tlist *C.List, scan_clauses *C.List, outer_plan *C.Plan,
) *C.ForeignScan {
	log.Printf("%p GetForeignPlan", root)

	id := C.lcons(nil, nil)
	state := initialized.execs.Load(root)
	initialized.execs.Delete(root)
	initialized.execs.Store(unsafe.Pointer(id), state)

	return C.make_foreignscan(
		tlist,
		C.goExtractBareClauses(scan_clauses),
		baserel.relid,
		nil, // TODO parameters
		id,
		nil, // TODO custom tlist
		nil, // TODO remote quals
		outer_plan)
}

// https://www.postgresql.org/docs/current/static/fdw-callbacks.html#FDW-CALLBACKS-SCAN
//
// Executing

//export goBeginForeignScan
func goBeginForeignScan(node *C.ForeignScanState, eflags C.int) {
	log.Printf("%p BeginForeignScan", node)

	fs := (*C.ForeignScan)(unsafe.Pointer(node.ss.ps.plan))
	node.fdw_state = unsafe.Pointer(fs.fdw_private)

	state := initialized.execs.Load(node.fdw_state)

	// TODO call the handler
	if eflags&C.EXEC_FLAG_EXPLAIN_ONLY != 0 {
		return
	}

	descriptor := node.ss.ss_ScanTupleSlot.tts_tupleDescriptor

	state.attributes = make([]Attribute, descriptor.natts)
	state.attinmeta = C.TupleDescGetAttInMetadata(descriptor)
	state.iterator = state.scanPath.Begin()

	for i, _ := range state.attributes {
		state.attributes[i] = attribute{index: C.int(i), state: state}
	}
}

//export goIterateForeignScan
func goIterateForeignScan(node *C.ForeignScanState) *C.TupleTableSlot {
	log.Printf("%p IterateForeignScan", node)

	state := initialized.execs.Load(node.fdw_state)
	state.slot = node.ss.ss_ScanTupleSlot

	C.goClearTupleSlot(state.slot)

	if state.iterator.HasNext() {
		state.iterator.Next(state.attributes)
		C.ExecStoreVirtualTuple(state.slot)
	}

	return state.slot

	// ExecStoreTuple(something, slot, InvalidBuffer, false);

	// rel = node->ss.ss_currentRelation;
	// attinmeta = TupleDescGetAttInMetadata(rel->rd_att);
	// natts = rel->rd_att->natts;
	// values = (char **) palloc(sizeof(char *) * natts);
	// for(i = 0; i < natts; i++ ){
	//   values[i] = "Hello,World";
	// }
	// tuple = BuildTupleFromCStrings(attinmeta, values);
	// ExecStoreTuple(something, slot, InvalidBuffer, true);
}

//export goReScanForeignScan
func goReScanForeignScan(node *C.ForeignScanState) {
	log.Printf("%p ReScanForeignScan", node)

	// TODO state.iterator rescan
}

//export goEndForeignScan
func goEndForeignScan(node *C.ForeignScanState) {
	log.Printf("%p EndForeignScan", node)

	state := initialized.execs.Load(node.fdw_state)
	initialized.execs.Delete(node.fdw_state)

	if state.iterator != nil {
		state.iterator.Close() // TODO error
	}
	state.scanPath.Close()
}
