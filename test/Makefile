
MODULE_big = go_test_fdw
OBJS = go_test_fdw.o

SHLIB_LINK = go_test_fdw.a

EXTENSION = go_test_fdw
DATA = go_test_fdw--1.0.sql

EXTRA_CLEAN = go_test_fdw.a go_test_fdw.h

REGRESS = --load-extension=$(MODULE_big) options errors joins

PG_CONFIG = pg_config
PGXS := $(shell $(PG_CONFIG) --pgxs)
include $(PGXS)

go_test_fdw.o: go_test_fdw.a
go_test_fdw.a: *.go
	CGO_CFLAGS='-I$(shell $(PG_CONFIG) --includedir-server)' go build -buildmode=c-archive -o $@ *.go
