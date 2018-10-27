package fdw

// https://www.postgresql.org/docs/current/static/fdw-callbacks.html#FDW-CALLBACKS-SCAN
//
// Executing

/*
#include "go_fdw.h"
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
	return TupleDescAttr(slot->tts_tupleDescriptor, i)->atttypid;
}

static ErrorData *
goSlotSetText(TupleTableSlot *slot, AttInMetadata *attinmeta, int i, char *value)
{
	volatile MemoryContext context = CurrentMemoryContext;
	ErrorData *edata = NULL;

	PG_TRY();
	{
		slot->tts_isnull[i] = (value == NULL) ? true : false;
		slot->tts_values[i] = InputFunctionCall(
			&attinmeta->attinfuncs[i], value, attinmeta->attioparams[i], attinmeta->atttypmods[i]);
	}
	PG_CATCH();
	{
		MemoryContextSwitchTo(context);
		edata = CopyErrorData();
		FlushErrorState();
	}
	PG_END_TRY();

	return edata;
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

func (a attribute) SetText(input []byte) error {
	buffer := make([]byte, len(input)+1)
	buffer[len(input)] = 0
	copy(buffer, input)
	return a.SetText0(buffer)
}

func (a attribute) SetText0(input []byte) error {
	return goErrorData(C.goSlotSetText(a.state.slot, a.state.attinmeta, a.index, (*C.char)(unsafe.Pointer(&input[0]))))
}

func (a attribute) SetString(input string) error {
	buffer := make([]byte, len(input)+1)
	buffer[len(input)] = 0
	copy(buffer, input)
	return a.SetText0(buffer)
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
func goBeginForeignScan(node *C.ForeignScanState, eflags C.int) *C.ErrorData {
	log.Printf("BeginForeignScan   (%p) Relation: %v", node, node.ss.ss_currentRelation.rd_id)

	// C.CurrentMemoryContext.name == "ExecutorState"

	if len(paths) > 0 {
		panic("Expected paths to be empty")
	}

	fs := (*C.ForeignScan)(unsafe.Pointer(node.ss.ps.plan))
	id := *(**C.RelOptInfo)(unsafe.Pointer(&fs.fdw_private.head.data[0]))
	state := &scan{plan: plans[id]}
	delete(plans, id)

	// TODO call the handler
	if eflags&C.EXEC_FLAG_EXPLAIN_ONLY != 0 {
		return nil
	}

	descriptor := node.ss.ss_ScanTupleSlot.tts_tupleDescriptor

	var err error
	state.attributes = make([]Attribute, descriptor.natts)
	state.attinmeta = C.TupleDescGetAttInMetadata(descriptor)
	state.iterator, err = state.plan.Begin()

	if err != nil {
		// TODO close()?
		initialized.handler = nil
		return pgErrorData(err)
	}

	for i, _ := range state.attributes {
		state.attributes[i] = attribute{index: C.int(i), state: state}
	}

	scans[node] = state

	return nil
}

//export goIterateForeignScan
func goIterateForeignScan(node *C.ForeignScanState) C.struct_goIterateForeignScanResult {
	log.Printf("IterateForeignScan (%p)", node)

	// C.CurrentMemoryContext.parent.name == "ExecutorState"
	// C.CurrentMemoryContext.name == "ExprContext"

	if len(plans) > 0 {
		panic("Expected plans to be empty")
	}

	state := scans[node]
	state.slot = node.ss.ss_ScanTupleSlot

	C.goClearTupleSlot(state.slot)

	if ok, err := state.iterator.Next(state.attributes); err != nil {
		delete(scans, node)
		// TODO close()?
		initialized.handler = nil
		return C.struct_goIterateForeignScanResult{edata: pgErrorData(err)}
	} else if ok {
		C.ExecStoreVirtualTuple(state.slot)
	}

	return C.struct_goIterateForeignScanResult{slot: state.slot}

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
func goEndForeignScan(node *C.ForeignScanState) *C.ErrorData {
	log.Printf("EndForeignScan     (%p) Relation: %v", node, node.ss.ss_currentRelation.rd_id)

	// C.CurrentMemoryContext.name == "ExecutorState"

	state := scans[node]
	delete(scans, node)

	var edata *C.ErrorData
	if state.iterator != nil {
		if err := state.iterator.Close(); err != nil {
			edata = pgErrorData(err)
		}
	}
	if err := state.plan.Close(); edata == nil && err != nil {
		edata = pgErrorData(err)
	}

	if len(scans) == 0 {
		// TODO close()?
		initialized.handler = nil
	}

	return edata
}
