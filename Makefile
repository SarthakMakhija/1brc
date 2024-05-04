build:
	go build -v ./...

test:
	go clean -testcache && go test -v ./...

build_and_test: build test

executable:
	go build -o cmd/main cmd/main.go

build_executable: build executable

clean:
	rm -rf cmd/main && rm -rf output.txt

all: clean build_and_test build_executable

run_10M:
	./run.sh

clean_run_10M: all run_10M

