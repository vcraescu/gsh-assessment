.PHONY: test test-local build clean sls-deploy gomodgen start start-local docker-build

test:
	docker build -t gs-assessment --no-cache --progress plain --target test-stage  .

test-local:
	go test ./...

start: docker-build
	docker run --rm -p 127.0.0.1:3000:3000 gs-assessment

start-local:
	go run ./cmd/app/main.go

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/lambda cmd/lambda/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/app cmd/app/main.go

clean:
	rm -rf ./bin

sls-deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

docker-build:
	docker build -t gs-assessment --progress plain --no-cache --target run-stage .
