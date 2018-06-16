-- vim: set expandtab shiftwidth=2 syntax=pgsql tabstop=2 :

CREATE SERVER none FOREIGN DATA WRAPPER go_test_fdw;

CREATE FOREIGN TABLE one (id integer, col1 text)
SERVER none OPTIONS (test 'fdw_join');

CREATE FOREIGN TABLE two (id integer, col2 text)
SERVER none OPTIONS (test 'fdw_join');

--
DO $$ BEGIN RAISE INFO 'local join'; END $$;
--
SELECT * FROM one JOIN (VALUES (1, 'a'), (2, 'b'), (3, 'c')) x (id, value) USING (id) ORDER BY id;

--
DO $$ BEGIN RAISE INFO 'fdw self-join'; END $$;
--
SELECT * FROM one a JOIN one b USING (id) ORDER BY id;

--
DO $$ BEGIN RAISE INFO 'fdw inner join'; END $$;
--
SELECT * FROM one JOIN two USING (id) ORDER BY id;

--
DO $$ BEGIN RAISE INFO 'fdw left join'; END $$;
--
SELECT * FROM one LEFT JOIN two USING (id) ORDER BY id;

--
DO $$ BEGIN RAISE INFO 'fdw right join'; END $$;
--
SELECT * FROM one RIGHT JOIN two USING (id) ORDER BY id;
