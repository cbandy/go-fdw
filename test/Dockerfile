ARG GO_VERSION
FROM golang:${GO_VERSION}-stretch
ARG GO_VERSION
ARG PG_MAJOR

RUN export DEBIAN_FRONTEND=noninteractive \
    && curl -sS 'https://www.postgresql.org/media/keys/ACCC4CF8.asc' | apt-key add - \
    && echo deb http://apt.postgresql.org/pub/repos/apt/ stretch-pgdg main "$PG_MAJOR" > /etc/apt/sources.list.d/postgresql.list \
    && apt-get update \
    && apt-get install --no-install-recommends --yes "postgresql-common" \
    && echo 'create_main_cluster = false' >> /etc/postgresql-common/createcluster.conf \
    && apt-get install --no-install-recommends --yes "postgresql-${PG_MAJOR}" "postgresql-server-dev-${PG_MAJOR}"
