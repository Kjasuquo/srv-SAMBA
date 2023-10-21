export GOPRIVATE=github.com/SAMBA-Research
export SERVICE_NAME=$(shell basename ${PWD})

build:
	cd src/cmd && go build -buildvcs=false -o service

tidy:
	cd src && go mod tidy

run:
	cd src && go run cmd/main.go