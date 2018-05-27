/* vim: set noexpandtab autoindent cindent tabstop=4 shiftwidth=4 cinoptions="(0,t0": */

#include "go_test_fdw.h"

PG_MODULE_MAGIC;

extern Datum go_test_fdw_handler(PG_FUNCTION_ARGS);
extern Datum go_test_fdw_validator(PG_FUNCTION_ARGS);

PG_FUNCTION_INFO_V1(go_test_fdw_handler);
PG_FUNCTION_INFO_V1(go_test_fdw_validator);

Datum
go_test_fdw_handler(PG_FUNCTION_ARGS)
{
	FdwRoutine *fr = makeNode(FdwRoutine);
	goTestInitialize(fr);

	// https://www.postgresql.org/docs/current/static/fdw-callbacks.html#FDW-CALLBACKS-JOIN-SCAN
	//fr->GetForeignJoinPaths;

	// https://www.postgresql.org/docs/current/static/fdw-callbacks.html#FDW-CALLBACKS-EXPLAIN
	//fr->ExplainForeignScan;

	// https://www.postgresql.org/docs/current/static/fdw-callbacks.html#FDW-CALLBACKS-ANALYZE
	//fr->AnalyzeForeignTable;

	// https://www.postgresql.org/docs/current/static/fdw-callbacks.html#FDW-CALLBACKS-IMPORT
	//fr->ImportForeignSchema;

	// https://www.postgresql.org/docs/current/static/fdw-callbacks.html#FDW-CALLBACKS-UPPER-PLANNING
	//fr->GetForeignUpperPaths;

	PG_RETURN_POINTER(fr);
}

Datum
go_test_fdw_validator(PG_FUNCTION_ARGS)
{
	/* no-op */
	PG_RETURN_VOID();
}
