# Define the output directories
BUILD_DIR = build
LINUX_DIR = $(BUILD_DIR)/linux
WINDOWS_DIR = $(BUILD_DIR)/windows

# Define the source file
SRC_FILE = main.go

# Define the output files
LINUX_OUTPUT = $(LINUX_DIR)/socket
WINDOWS_OUTPUT = $(WINDOWS_DIR)/socket.exe

# Default target
all: build-linux build-windows zip

# Build for Linux
build-linux: $(LINUX_OUTPUT)

$(LINUX_DIR)/socket: $(SRC_FILE)
	mkdir -p $(LINUX_DIR)
	GOOS=linux GOARCH=amd64 go build -o $@ $<

# Build for Windows
build-windows: $(WINDOWS_OUTPUT)

$(WINDOWS_DIR)/socket.exe: $(SRC_FILE)
	mkdir -p $(WINDOWS_DIR)
	GOOS=windows GOARCH=amd64 go build -o $@ $<

# Zip the files
zip: $(LINUX_OUTPUT) $(WINDOWS_OUTPUT)
	zip -j $(BUILD_DIR)/sockets.zip $(LINUX_OUTPUT) $(WINDOWS_OUTPUT)

# Clean the build directory
clean:
	rm -rf $(BUILD_DIR)

.PHONY: all build-linux build-windows zip clean
