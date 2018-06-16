-- vim: set expandtab shiftwidth=2 syntax=pgsql tabstop=2 :

--
DO $$ BEGIN RAISE INFO 'no server options'; END $$;
--
CREATE SERVER none FOREIGN DATA WRAPPER go_test_fdw;
CREATE FOREIGN TABLE server_options (option text, value text)
SERVER none;

-- expect empty set
SELECT * FROM server_options ORDER BY option;

--
DO $$ BEGIN RAISE INFO 'some server options'; END $$;
--
DROP FOREIGN TABLE server_options;
CREATE SERVER few FOREIGN DATA WRAPPER go_test_fdw
OPTIONS (uno 'dos', tres 'quatro');
CREATE FOREIGN TABLE server_options (option text, value text)
SERVER few;

-- expect non-empty set
SELECT * FROM server_options ORDER BY option;

--
DO $$ BEGIN RAISE INFO 'no table options'; END $$;
--
CREATE FOREIGN TABLE table_options (option text, value text)
SERVER few;

-- expect empty set
SELECT * FROM table_options ORDER BY option;

--
DO $$ BEGIN RAISE INFO 'some table options'; END $$;
--
DROP FOREIGN TABLE table_options;
CREATE FOREIGN TABLE table_options (option text, value text)
SERVER few
OPTIONS (ichi 'ni', san 'shi');

-- expect non-empty set
SELECT * FROM table_options ORDER BY option;

--
DO $$ BEGIN RAISE INFO 'no user mapping'; END $$;
--
CREATE FOREIGN TABLE user_options (option text, value text)
SERVER few;

-- expect error
SELECT * FROM user_options ORDER BY option;

--
DO $$ BEGIN RAISE INFO 'empty user mapping'; END $$;
--
CREATE USER MAPPING FOR PUBLIC
SERVER few;

-- expect empty set
SELECT * FROM user_options ORDER BY option;

--
DO $$ BEGIN RAISE INFO 'some user mapping'; END $$;
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
