-- vim: set expandtab shiftwidth=2 syntax=pgsql tabstop=2 :

CREATE EXTENSION go_test_fdw;

CREATE SERVER gts FOREIGN DATA WRAPPER go_test_fdw
OPTIONS ( test 'options', uno 'dos', tres 'quatro' );

CREATE FOREIGN TABLE table_options ( option text, value text )
SERVER gts
OPTIONS ( tres 'replaced' );

-- expect combined table and server options
SELECT * FROM table_options ORDER BY option;
