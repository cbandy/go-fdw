#ifndef GO_FDW_H
#define GO_FDW_H

#include "postgres.h"
#include "foreign/fdwapi.h"
#include "utils/elog.h"

ErrorData *
goGetForeignRelSize(
	PlannerInfo *root, RelOptInfo *baserel, Oid foreigntableid);

void
goGetForeignPaths(
	PlannerInfo *root, RelOptInfo *baserel, Oid foreigntableid);

ForeignScan *
goGetForeignPlan(
	PlannerInfo *root, RelOptInfo *baserel, Oid foreigntableid,
	ForeignPath *best_path, List *tlist, List *scan_clauses, Plan *outer_plan);

ErrorData *
goBeginForeignScan(ForeignScanState *node, int eflags);

struct goIterateForeignScanResult {
	TupleTableSlot *slot;
	ErrorData *edata;
};

struct goIterateForeignScanResult
goIterateForeignScan(ForeignScanState *node);

void
goReScanForeignScan(ForeignScanState *node);

ErrorData *
goEndForeignScan(ForeignScanState *node);

#endif
