.SILENT:

build:
	go build -o bin/app cmd/app/main.go

run: build
	./bin/app