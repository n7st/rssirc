build:
	go build -o rssirc cmd/rssirc/main.go

test:
	go test -v ./...