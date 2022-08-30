.PHONY: test
test: # Run unit tests
	go test ./...

.PHONY: build
build: # Build package
	mkdir -p build/
	go build -o build/

.PHONY: clean
clean: # Clean build folder
	rm -rf build/
