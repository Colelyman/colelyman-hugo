
all: get build

get:
	git submodule update --init --recursive

build:
	go get ./functions/...
	go build -o functions/micropub ./functions/...
