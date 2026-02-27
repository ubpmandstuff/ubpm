binary_name := "ubpm"
version := "0.1.0+beta"
destdir := "/usr/local/bin"

run:
	go mod tidy
	go run .

build: clean
	mkdir -p "build/{{version}}"
	go mod tidy
	go build -o "build/{{version}}/{{binary_name}}"

build-all: clean build-win build-osx build-lin

build-win:
	mkdir -p "build/{{version}}"
	go mod tidy
	GOARCH=amd64 GOOS=windows go build -o "build/{{version}}/{{binary_name}}-windows-amd64.exe"

build-osx:
	mkdir -p "build/{{version}}"
	go mod tidy
	GOARCH=arm64 GOOS=darwin go build -o "build/{{version}}/{{binary_name}}-osx-arm64"

build-lin:
	mkdir -p "build/{{version}}"
	go mod tidy
	GOARCH=amd64 GOOS=linux go build -o "build/{{version}}/{{binary_name}}-linux-amd64"

# install: build
# 	mkdir -p {{destdir}}
# 	cp build/{{version}}/{{binary_name}} {{destdir}}/{{binary_name}}

clean:
	@echo "cleaning builds for current version"
	rm -rf build/{{version}}/*

clean-all:
	@echo "cleaning all builds"
	rm -rf build/*
