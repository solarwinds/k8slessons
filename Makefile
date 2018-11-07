PHONY: build run

app_name := "foobar"

build:
	@go build -i -o $(app_name)

run: build
run:
	@./$(app_name)

