.PHONY: build test

test:
	docker-compose run --rm app gotestsum --format short-verbose --packages="./..." -- -p 1

build:
	docker-compose run --rm app bash -c 'rm -rf build && go build -o build/server main.go'
