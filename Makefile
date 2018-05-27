
.PHONY: check
check:
	docker build --file test/Dockerfile --label 'check=go-fdw' .

.PHONY: clean
clean:
	docker image prune --filter 'label=check=go-fdw'
