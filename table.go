package fdw

/*
#include "postgres.h"
#include "foreign/fdwapi.h"

#include "commands/defrem.h"
#include "foreign/foreign.h"
#include "miscadmin.h"
#include "utils/rel.h"

static inline char *
goRelationName(Relation relation)
{
	return RelationGetRelationName(relation);
}
*/
import "C"
import "unsafe"

type table struct{ oid C.Oid }
type tableOptions struct{ ft *C.ForeignTable }
type tableRelation struct{ rel C.Relation }

func (t table) Oid() uint { return uint(t.oid) }

func (t table) Options(f func(Options)) {
	f(tableOptions{C.GetForeignTable(t.oid)})
}

func (t table) Relation(f func(Relation)) {
	rel := C.RelationIdGetRelation(t.oid)
	f(tableRelation{rel})
	C.RelationClose(rel)
}

func (tableOptions) makeMap(defElems *C.List) map[string]string {
	result := make(map[string]string)

	if defElems != nil {
		for lc := defElems.head; lc != nil; lc = lc.next {
			de := *(**C.DefElem)(unsafe.Pointer(&lc.data[0]))
			result[C.GoString(de.defname)] = C.GoString(C.defGetString(de))
		}
	}

	return result
}

func (o tableOptions) Server() map[string]string {
	return o.makeMap(C.GetForeignServer(o.ft.serverid).options)
}

func (o tableOptions) Table() map[string]string {
	return o.makeMap(o.ft.options)
}

func (o tableOptions) User() map[string]string {
	return o.makeMap(C.GetUserMapping(C.GetUserId(), o.ft.serverid).options)
}

func (r tableRelation) Name() string {
	return C.GoString(C.goRelationName(r.rel))
}
