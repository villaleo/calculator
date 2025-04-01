# Build the server into a standalone binary.
server:
	go build -o build/calc .

# Clean the entire build directory.
clean:
	rm -rf build/*
