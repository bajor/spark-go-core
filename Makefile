all: test run

test:
	go test -count=1 ./operations/...
	go test -count=1 ./rdd/...

run:
	go run main.go 