ARG GO_VERSION=1.10
FROM golang:${GO_VERSION}-stretch
ARG GO_VERSION
ARG PG_MAJOR=10

RUN export DEBIAN_FRONTEND=noninteractive \
    && curl -sS 'https://www.postgresql.org/media/keys/ACCC4CF8.asc' | apt-key add - \
    && echo deb http://apt.postgresql.org/pub/repos/apt/ stretch-pgdg main "$PG_MAJOR" > /etc/apt/sources.list.d/postgresql.list \
    && apt-get update \
    && apt-get install --no-install-recommends --yes "postgresql-common" \
    && echo 'create_main_cluster = false' >> /etc/postgresql-common/createcluster.conf \
    && apt-get install --no-install-recommends --yes "postgresql-${PG_MAJOR}" "postgresql-server-dev-${PG_MAJOR}"

COPY .  /go/src/github.com/cbandy/go-fdw
WORKDIR /go/src/github.com/cbandy/go-fdw/test

RUN make clean all install
RUN chown -R postgres: . && su -c 'make REGRESS_OPTS="--temp-instance=$(mktemp -d)" installcheck' postgres \
    || { rc=$? ; cat regression.diffs log/postmaster.log ; exit $rc ; }
