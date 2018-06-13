-- vim: set expandtab shiftwidth=2 syntax=pgsql tabstop=2 :

--
-- no server options
--
CREATE SERVER none FOREIGN DATA WRAPPER go_test_fdw;
CREATE FOREIGN TABLE server_options (option text, value text)
SERVER none;

-- expect empty set
SELECT * FROM server_options ORDER BY option;

--
-- some server options
--
DROP FOREIGN TABLE server_options;
CREATE SERVER few FOREIGN DATA WRAPPER go_test_fdw
OPTIONS (uno 'dos', tres 'quatro');
CREATE FOREIGN TABLE server_options (option text, value text)
SERVER few;

-- expect non-empty set
SELECT * FROM server_options ORDER BY option;

--
-- no table options
--
CREATE FOREIGN TABLE table_options (option text, value text)
SERVER few;

-- expect empty set
SELECT * FROM table_options ORDER BY option;

--
-- some table options
--
DROP FOREIGN TABLE table_options;
CREATE FOREIGN TABLE table_options (option text, value text)
SERVER few
OPTIONS (ichi 'ni', san 'shi');

-- expect non-empty set
SELECT * FROM table_options ORDER BY option;

--
-- no user mapping
--
CREATE FOREIGN TABLE user_options (option text, value text)
SERVER few;

-- expect error
SELECT * FROM user_options ORDER BY option;

--
-- empty user mapping
--
CREATE USER MAPPING FOR PUBLIC
SERVER few;

-- expect empty set
SELECT * FROM user_options ORDER BY option;

--
-- some user mapping
--
CREATE USER MAPPING FOR postgres
SERVER few
OPTIONS (eins 'zwei', drei 'vier');

-- expect non-empty set
SELECT * FROM user_options ORDER BY option;

--
-- cleanup
--
DROP SERVER few, none CASCADE;
