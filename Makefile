.PHONY: validate

SOURCE_FILES  := ./pkg/...

validate:
	go fmt ./...

clean:
	go clean
	rm -rf bin

build: clean
	mkdir bin
	#GOOS=linux GOARCH=amd64 go build -o bin/ ./...
	go build -o bin/ ./...

coverage: test
	go tool cover -func coverage.out

test:
	go test -cover -coverprofile=coverage.out ${SOURCE_FILES}

view-coverage: test
	go tool cover -html=coverage.out

run:
	./bin/app

prebuild-checks: validate test coverage