
IMAGE ?= go-fdw:check
GO_VERSION ?= 1.11
PG_MAJOR ?= 10

.PHONY: build-check
build-check:
	docker build --label 'check=go-fdw' --build-arg 'GO_VERSION=$(GO_VERSION)' --build-arg 'PG_MAJOR=$(PG_MAJOR)' --tag $(IMAGE) test

.PHONY: check
check: build-check
	docker run   --label 'check=go-fdw' --rm --mount "target=/mnt,source=$$(pwd),type=bind" \
		$(IMAGE) /mnt/test/docker-check.sh

.PHONY: clean
clean:
	docker image prune --filter 'label=check=go-fdw'
