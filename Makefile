build:
	go build -v ./...

test:
	go clean -testcache && go test -v ./...

build_and_test: build test

executable:
	go build -o cmd/main cmd/main.go

build_executable: build executable

clean:
	rm -rf cmd/main && rm -rf output.txt && rm -rf profile.out

all: clean build_and_test build_executable

run_10M:
	./run.sh

clean_run_10M: all run_10M

profile_10M:
	cmd/main -f ./fixture/10M_weather_stations.csv -cpuprofile profile.out

clean_profile_10M: all profile_10M
