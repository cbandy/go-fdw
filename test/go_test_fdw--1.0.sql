-- vim: set expandtab shiftwidth=2 syntax=pgsql tabstop=2 :

CREATE FUNCTION go_test_fdw_benchmark_begin()
  RETURNS void
  LANGUAGE C STRICT
AS 'MODULE_PATHNAME';

CREATE FUNCTION go_test_fdw_benchmark_end()
  RETURNS text
  LANGUAGE C STRICT
AS 'MODULE_PATHNAME';

-- FDW

CREATE FUNCTION go_test_fdw_handler()
  RETURNS fdw_handler
  LANGUAGE C STRICT
AS 'MODULE_PATHNAME';

CREATE FUNCTION go_test_fdw_validator(text[], oid)
  RETURNS void
  LANGUAGE C STRICT
AS 'MODULE_PATHNAME';

CREATE FOREIGN DATA WRAPPER go_test_fdw
  HANDLER go_test_fdw_handler
  VALIDATOR go_test_fdw_validator;
