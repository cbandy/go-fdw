
IMAGE ?= go-fdw:check

.PHONY: check
check:
	docker build --label 'check=go-fdw' --tag $(IMAGE) test
	docker run   --label 'check=go-fdw' --rm --mount "target=/mnt,source=$$(pwd),type=bind" \
		$(IMAGE) /mnt/test/docker-check.sh

.PHONY: clean
clean:
	docker image prune --filter 'label=check=go-fdw'
