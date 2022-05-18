.PHONY : all

ts=$(shell date +%F_%H-%M-%S)

all : static
	cp ./static index.tmpl.html *.js *.css ./docker/
	docker build docker -t gopf:${ts}
	docker tag gopf:${ts} gopf:latest

static : main.go
	CGO_ENABLED=0 go build -o $@ $^
