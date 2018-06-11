
IMAGE ?= go-fdw:check

.PHONY: check
check:
	docker build --file test/Dockerfile --label 'check=go-fdw' --tag $(IMAGE) .
	docker run --rm --interactive --tty --label 'check=go-fdw' --mount "target=/mnt,source=$$(pwd),type=bind" \
		$(IMAGE) /mnt/test/docker-check.sh

.PHONY: clean
clean:
	docker image prune --filter 'label=check=go-fdw'
