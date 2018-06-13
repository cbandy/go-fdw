-- vim: set expandtab shiftwidth=2 syntax=pgsql tabstop=2 :

CREATE SERVER none FOREIGN DATA WRAPPER go_test_fdw;
CREATE FOREIGN TABLE table_options (option text, value text)
SERVER none OPTIONS (a 'b', c 'd', e 'f');

--
-- local join
--
SELECT * FROM table_options JOIN (VALUES ('a', '1'), ('b', '2'), ('c', '3')) x (option, number) USING (option)
ORDER BY option;
