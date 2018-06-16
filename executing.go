package fdw

// https://www.postgresql.org/docs/current/static/fdw-callbacks.html#FDW-CALLBACKS-SCAN
//
// Executing

/*
#include "postgres.h"
#include "foreign/fdwapi.h"

#include "funcapi.h"

static inline void
goClearTupleSlot(TupleTableSlot *slot)
{
	ExecClearTuple(slot);
	memset(slot->tts_values, 0, sizeof(Datum) * slot->tts_tupleDescriptor->natts);
	memset(slot->tts_isnull, true, sizeof(bool) * slot->tts_tupleDescriptor->natts);
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
*/
import "C"
import (
	"log"
	"unsafe"
)

type attribute struct {
	index C.int
	state *scan
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

func (a attribute) SetString(input string) {
	buffer := make([]byte, len(input)+1)
	buffer[len(input)] = 0
	copy(buffer, input)
	a.SetText0(buffer)
}

func (a attribute) TypeOid() uint {
	return uint(C.goSlotGetTypeOid(a.state.slot, a.index))
}

type scan struct {
	plan ScanPath // TODO generic plan

	iterator Iterator

	attributes []Attribute
	attinmeta  *C.AttInMetadata
	slot       *C.TupleTableSlot
}

var scans = map[*C.ForeignScanState]*scan{}

//export goBeginForeignScan
func goBeginForeignScan(node *C.ForeignScanState, eflags C.int) {
	log.Printf("BeginForeignScan   (%p) Relation: %v", node, node.ss.ss_currentRelation.rd_id)

	// C.CurrentMemoryContext.name == "ExecutorState"

	if len(paths) > 0 {
		panic("Expected paths to be empty")
	}

	fs := (*C.ForeignScan)(unsafe.Pointer(node.ss.ps.plan))
	state := &scan{plan: plans[fs]}
	delete(plans, fs)

	// TODO call the handler
	if eflags&C.EXEC_FLAG_EXPLAIN_ONLY != 0 {
		return
	}

	descriptor := node.ss.ss_ScanTupleSlot.tts_tupleDescriptor

	state.attributes = make([]Attribute, descriptor.natts)
	state.attinmeta = C.TupleDescGetAttInMetadata(descriptor)
	state.iterator = state.plan.Begin()

	for i, _ := range state.attributes {
		state.attributes[i] = attribute{index: C.int(i), state: state}
	}

	scans[node] = state
}

//export goIterateForeignScan
func goIterateForeignScan(node *C.ForeignScanState) *C.TupleTableSlot {
	log.Printf("IterateForeignScan (%p)", node)

	// C.CurrentMemoryContext.parent.name == "ExecutorState"
	// C.CurrentMemoryContext.name == "ExprContext"

	if len(plans) > 0 {
		panic("Expected plans to be empty")
	}

	state := scans[node]
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
	log.Printf("EndForeignScan     (%p) Relation: %v", node, node.ss.ss_currentRelation.rd_id)

	// C.CurrentMemoryContext.name == "ExecutorState"

	state := scans[node]
	delete(scans, node)

	if state.iterator != nil {
		state.iterator.Close() // TODO error
	}
	state.plan.Close()

	if len(scans) == 0 {
		// TODO close()?
		initialized.handler = nil
	}
}
