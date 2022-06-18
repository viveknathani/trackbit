build:
	go build -o ./bin/trackbit main.go

test:
	go test -v ./...

run:
	export PORT=8080 && ./bin/trackbit

