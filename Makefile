.PHONY : all

ts=$(shell date +%F_%H-%M-%S)

all : static
	docker build . -t gopf:${ts}
	docker tag gopf:${ts} gopf:latest

static : main.go
	CGO_ENABLED=0 go build -o $@ $^
