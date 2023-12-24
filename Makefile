GIT_COMMIT=$(shell git rev-parse --short HEAD)

dist: clean site
	./site build --include-unpublished

site:
	go build

.PHONY: clean
clean:
	- rm site
	- rm -rf dist

.PHONY: debug
debug:
	@python3 -m http.server -d dist

.PHONY: docker
docker:
	@docker build -t $(GIT_COMMIT) . && \
		docker run --rm -it -p 8080:8080 $(GIT_COMMIT)
