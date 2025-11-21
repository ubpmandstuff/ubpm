binary_name := "ubpm"
version := "0.0.2+alpha"

run:
	go mod tidy
	go run .

build: clean
	mkdir -p "build/{{version}}"
	go mod tidy
	go build -o "build/{{version}}/{{binary_name}}"

clean:
	@echo "cleaning builds for current version"
	rm -rf build/{{version}}/*

clean-all:
	@echo "cleaning all builds"
	rm -rf build/*
