.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./perf
	rm -rf ./perf-mac
	
build:
	GOOS=linux GOARCH=amd64 go build -o perf .
	go build -o perf-mac .

run:
	./perf-mac -debug false

