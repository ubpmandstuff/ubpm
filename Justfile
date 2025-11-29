binary_name := "ubpm"
version := "0.0.4+alpha"
destdir := "/usr/local/bin"

run:
	go mod tidy
	go run .

build: clean
	mkdir -p "build/{{version}}"
	go mod tidy
	go build -o "build/{{version}}/{{binary_name}}"

install: build
	mkdir -p {{destdir}}
	cp build/{{version}}/{{binary_name}} {{destdir}}/{{binary_name}}

clean:
	@echo "cleaning builds for current version"
	rm -rf build/{{version}}/*

clean-all:
	@echo "cleaning all builds"
	rm -rf build/*
