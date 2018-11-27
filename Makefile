PHONY: build run

app_name := "dumbstore"
latest_git_sha := $(shell git rev-parse HEAD)
the_date := $(shell date +'%s')

web_app_name := $(app_name)-web

build:
	@go build -i -o $(app_name)

run: build
run:
	@./$(app_name)

run_compose:
	docker-compose up --build

docker_build:
	docker build -t $(app_name) .

bounce_web_app: docker_build
	kubectl patch deployment $(web_app_name) -p "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"date\":\"$(the_date)\"}}}}}"
