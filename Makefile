build:
	go build -o ./bin/trackbit main.go

test:
	go test -v ./...

run:
	./bin/trackbit

