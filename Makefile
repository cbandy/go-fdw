
IMAGE ?= go-fdw:check
GO_VERSION ?= 1.11
PG_MAJOR ?= 11

.PHONY: build-check
build-check:
	docker build \
		--build-arg 'GO_VERSION=$(GO_VERSION)' \
		--build-arg 'PG_MAJOR=$(PG_MAJOR)' \
		--label 'check=go-fdw' \
		--tag $(IMAGE) test

.PHONY: check
check: build-check
	docker run --rm \
		--label 'check=go-fdw' \
		--mount "target=/mnt,source=$$(pwd),type=bind" \
		$(IMAGE) /mnt/test/docker-check.sh

.PHONY: clean
clean:
	docker image prune --filter 'label=check=go-fdw'
