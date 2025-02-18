.PHONY: run init build docker-build docker-run test clean doc tidy

run:
	air

init:
	make tidy && make doc && make build

build:
	go build -o bin/app cmd/main.go

docker-build:
	docker build -t go-auth-service .

docker-run:
	docker-compose up

test:
	go test ./...

clean:
	rm -rf bin/ tmp/

doc:
	swag init --generalInfo cmd/main.go --output ./docs

tidy:
	go mod tidy
