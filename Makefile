.PHONY: build
build:
	go build -o static cmd/main.go

.PHONY: up
up:
	./static