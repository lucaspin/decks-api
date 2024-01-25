.PHONY: build test

test:
	docker-compose run --rm app gotestsum --format short-verbose --packages="./..." -- -p 1

build:
	rm -rf build && go build -o build/server main.go

server.start:
	$(MAKE) server.stop
	docker-compose run --rm app bash -c 'rm -rf build && go build -o build/server main.go'
	docker-compose up -d

server.stop:
	docker-compose down

server.logs:
	docker-compose logs app
