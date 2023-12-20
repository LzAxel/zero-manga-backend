.SILENT:

build:
	go build -o bin/app cmd/app/main.go

run: build
	./bin/app

swag:
	swag init --parseInternal -d cmd/app/,internal/handler/http/ -o docs