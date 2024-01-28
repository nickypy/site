dist: clean site
	./site build --include-unpublished

site:
	go build

.PHONY: clean
clean:
	- rm site
	- rm -rf dist

.PHONY: serve
serve:
	@go run github.com/eliben/static-server@latest dist
