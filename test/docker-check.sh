#!/bin/bash
set -eu

mkdir -p /go/src/github.com/cbandy/go-fdw
cp -R -t /go/src/github.com/cbandy/go-fdw /mnt/*
cd       /go/src/github.com/cbandy/go-fdw/test

make clean all install

chown -R postgres: .
su -c 'make GODEBUG="cgocheck=2" PG_REGRESS_DIFF_OPTS="-U3" REGRESS_OPTS="--temp-config=postgresql.conf --temp-instance=$(mktemp -d)" installcheck' postgres ||
{ rc=$? ; cat regression.diffs log/postmaster.log ; exit $rc ; }
