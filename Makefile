# Makefile for building Fairway Bridge for different targets

# Variables
TARGET ?= local
BINARY_NAME = fairway-bridge
SRC = main.go
ASSETS_DIR = Assets
OUTPUT_DIR = bin

# Set environment variables based on TARGET
ifeq ($(TARGET), mac)
    GOOS = darwin
    GOARCH = amd64
    BIN_SUFFIX = _OSX
else ifeq ($(TARGET), windows)
    GOOS = windows
    GOARCH = amd64
    BIN_SUFFIX =
else ifeq ($(TARGET), rpi)
    GOOS = linux
    GOARCH = arm64
    BIN_SUFFIX = _RPI
else
    GOOS ?= $(shell go env GOOS)
    GOARCH ?= $(shell go env GOARCH)
    BIN_SUFFIX =
endif

BINARY_PATH = $(OUTPUT_DIR)/$(BINARY_NAME)$(BIN_SUFFIX)
ZIPFILE = $(OUTPUT_DIR)/$(TARGET)_$(BINARY_NAME).zip

# Default target: clean, build, and bundle
all: clean build bundle

# Build the binary for the specified target
build:
	mkdir -p $(OUTPUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_PATH) $(SRC)
	@echo "Built $(BINARY_PATH) for $(GOOS) $(GOARCH)"

# Bundle the binary and Assets into a zip file
bundle:
	mkdir -p bundle_temp
	cp -r $(ASSETS_DIR) bundle_temp/
	cp $(BINARY_PATH) bundle_temp/
	cd bundle_temp && zip -r ../$(ZIPFILE) .
	rm -rf bundle_temp
	@echo "Created bundle $(ZIPFILE)"

# Clean build artifacts
clean:
	rm -rf $(OUTPUT_DIR) bundle_temp
	@echo "Cleaned build artifacts"

.PHONY: all build bundle clean
