package fdw

/*
#include "postgres.h"
#include "foreign/fdwapi.h"

#include "commands/defrem.h"
#include "foreign/foreign.h"
#include "miscadmin.h"
#include "utils/rel.h"

struct goGetUserMappingResult {
	UserMapping *mapping;
	ErrorData *edata;
};

static inline struct goGetUserMappingResult
goGetUserMapping(Oid serverid)
{
	volatile MemoryContext context = CurrentMemoryContext;
	struct goGetUserMappingResult result = {NULL, NULL};

	PG_TRY();
	{
		result.mapping = GetUserMapping(GetUserId(), serverid); // foreign/foreign.h, miscadmin.h
	}
	PG_CATCH();
	{
		MemoryContextSwitchTo(context);
		result.edata = CopyErrorData();
		FlushErrorState();
	}
	PG_END_TRY();

	return result;
}

static inline char *
goRelationName(Relation relation)
{
	return RelationGetRelationName(relation); // utils/rel.h
}
*/
import "C"
import "unsafe"

type table struct{ oid C.Oid }
type tableOptions struct{ ft *C.ForeignTable }
type tableRelation struct{ rel C.Relation }

func (t table) Oid() uint { return uint(t.oid) }

func (t table) Options(f func(Options)) {
	f(tableOptions{C.GetForeignTable(t.oid)}) // foreign/foreign.h
}

func (t table) Relation(f func(Relation)) {
	rel := C.RelationIdGetRelation(t.oid) // utils/relcache.h
	f(tableRelation{rel})                 //
	C.RelationClose(rel)                  // utils/relcache.h
}

func (tableOptions) makeMap(defElems *C.List) map[string]string {
	result := make(map[string]string)

	if defElems != nil {
		for lc := defElems.head; lc != nil; lc = lc.next {
			de := *(**C.DefElem)(unsafe.Pointer(&lc.data[0]))               // nodes/parsenodes.h
			result[C.GoString(de.defname)] = C.GoString(C.defGetString(de)) // commands/defrem.h
		}
	}

	return result
}

func (o tableOptions) Server() map[string]string {
	return o.makeMap(C.GetForeignServer(o.ft.serverid).options) // foreign/foreign.h
}

func (o tableOptions) Table() map[string]string {
	return o.makeMap(o.ft.options)
}

func (o tableOptions) User() (map[string]string, error) {
	options := (*C.List)(nil)
	result := C.goGetUserMapping(o.ft.serverid)

	if result.mapping != nil {
		options = result.mapping.options
	}

	return o.makeMap(options), goErrorData(result.edata)
}

func (r tableRelation) Name() string {
	return C.GoString(C.goRelationName(r.rel))
}
