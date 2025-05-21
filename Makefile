# Makefile for BFS and DFS Simulator

# Go commands
GO := go
BUILD := $(GO) build
RUN := $(GO) run
TIDY := $(GO) mod tidy

# Paths
MAIN_DIR := ./cmd/simulator
OUTPUT := bfsdfs
SAVES_DIR := ./saves

.PHONY: all build run clean tidy

# Default target
all: build

# Create saves directory
dirs:
	mkdir -p $(SAVES_DIR)

# Build the application
build: dirs
	$(BUILD) -o $(OUTPUT) $(MAIN_DIR)

# Run the application
run: dirs
	$(RUN) $(MAIN_DIR)

# Clean build artifacts
clean:
	rm -f $(OUTPUT)

# Tidy go modules
tidy:
	$(TIDY)

# Build and run
dev: build
	./$(OUTPUT)
