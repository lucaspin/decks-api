.PHONY: build test

test:
	docker-compose run --rm app gotestsum --format short-verbose --packages="./..." -- -p 1

build:
	docker-compose run --rm app bash -c 'rm -rf build && go build -o build/server main.go'

start:
	$(MAKE) stop
	$(MAKE) build
	docker-compose up -d

stop:
	docker-compose down

logs:
	docker-compose logs app
