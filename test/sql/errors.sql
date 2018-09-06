-- vim: set expandtab shiftwidth=2 syntax=pgsql tabstop=2 :

\set SHOW_CONTEXT errors
\set VERBOSITY verbose

SET log_error_verbosity = VERBOSE;
CREATE SERVER errors FOREIGN DATA WRAPPER go_test_fdw OPTIONS (test 'errors');

--
CREATE FOREIGN TABLE scan_path (id text) SERVER errors;
--
SELECT * FROM scan_path;

--
CREATE FOREIGN TABLE estimate_scan (id text) SERVER errors;
--
SELECT * FROM estimate_scan;

--
CREATE FOREIGN TABLE begin_scan (id text) SERVER errors;
--
SELECT * FROM begin_scan;

--
CREATE FOREIGN TABLE during_scan (id text) SERVER errors;
--
SELECT * FROM during_scan;

--
CREATE FOREIGN TABLE bad_conversion (id integer) SERVER errors;
--
SELECT * FROM bad_conversion;

--
CREATE FOREIGN TABLE end_scan (id text) SERVER errors;
--
SELECT * FROM end_scan;

--
CREATE FOREIGN TABLE end_path (id text) SERVER errors;
--
SELECT * FROM end_path;

--
-- cleanup
--
\unset VERBOSITY
DROP SERVER errors CASCADE;
