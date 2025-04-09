# Set the name of the output binary
BINARY_NAME := server

# Set the output directory
BIN_DIR := bin

# Set the source file
SOURCE_FILE := cmd/server/main.go

# Default target to display help
all:
	@echo "Use 'make build' to build the server binary."

# Build target
build: $(BIN_DIR)/$(BINARY_NAME)

# Rule to build the binary
$(BIN_DIR)/$(BINARY_NAME): $(SOURCE_FILE)
	@mkdir -p $(BIN_DIR)                     # Create the bin directory if it doesn't exist
	go build -o $@ $<                       # Build the Go binary

# Clean up the built binary
clean:
	rm -f $(BIN_DIR)/$(BINARY_NAME)

.PHONY: all build clean