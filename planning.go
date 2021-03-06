package fdw

// https://www.postgresql.org/docs/current/static/fdw-callbacks.html#FDW-CALLBACKS-SCAN
//
// Planning

/*
#include "go_fdw.h"
#include "optimizer/pathnode.h"
#include "optimizer/planmain.h"
#include "optimizer/restrictinfo.h"
#include "utils/memutils.h"

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

static inline List *
goExtractBareClauses(List *restrictinfo_list)
{
	return extract_actual_clauses(restrictinfo_list, false);
}

static inline List *
goPrivatePlanId(RelOptInfo *baserel)
{
	List *result = list_make1_int(0);
	list_head(result)->data.ptr_value = baserel;
	return result;
}
*/
import "C"
import "log"

var paths = map[*C.RelOptInfo]*struct {
	scanCost ScanCostEstimate
	scanPath ScanPath
}{}

var plans = map[*C.RelOptInfo]ScanPath{} // TODO generic plan // TODO uintptr key would be better

//export goGetForeignRelSize
func goGetForeignRelSize(root *C.PlannerInfo, baserel *C.RelOptInfo, foreigntableid C.Oid) *C.ErrorData {
	log.Printf("GetForeignRelSize (%p, %p, %v)", root, baserel, foreigntableid)

	// C.CurrentMemoryContext.name == "MessageContext"
	// C.CurrentMemoryContext.name == "SPI Exec"

	if _, exists := paths[baserel]; exists {
		panic("RelOptInfo already seen")
	}

	if initialized.handler == nil {
		initialized.handler = initialized.fdw.New()
	}

	state := new(struct {
		scanCost ScanCostEstimate
		scanPath ScanPath
	})

	var err error
	state.scanPath, err = initialized.handler.Scan(table{foreigntableid})

	if err == nil {
		state.scanCost, err = state.scanPath.Estimate(ScanCostEstimate{
			Rows:  float64(baserel.rows),
			Width: int(baserel.reltarget.width),
		})
	}

	if err != nil {
		// TODO close()?
		initialized.handler = nil
		return pgErrorData(err)
	}

	baserel.rows = C.double(state.scanCost.Rows)
	baserel.reltarget.width = C.int(state.scanCost.Width)

	paths[baserel] = state

	return nil
}

//export goGetForeignPaths
func goGetForeignPaths(root *C.PlannerInfo, baserel *C.RelOptInfo, foreigntableid C.Oid) {
	log.Printf("GetForeignPaths   (%p, %p, %v)", root, baserel, foreigntableid)

	// C.CurrentMemoryContext.name == "MessageContext"
	// C.CurrentMemoryContext.name == "SPI Exec"

	state := paths[baserel]

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
	log.Printf("GetForeignPlan    (%p, %p, %v, %p)", root, baserel, foreigntableid, best_path)

	// C.CurrentMemoryContext.name == "MessageContext"
	// C.CurrentMemoryContext.name == "SPI Exec"

	state := paths[baserel]
	delete(paths, baserel)

	fs := C.make_foreignscan( // optimizer/planmain.h
		tlist,
		C.goExtractBareClauses(scan_clauses),
		baserel.relid,
		nil, // TODO parameters
		C.goPrivatePlanId(baserel),
		nil, // TODO custom tlist
		nil, // TODO remote quals
		outer_plan)

	plans[baserel] = state.scanPath
	return fs
}
