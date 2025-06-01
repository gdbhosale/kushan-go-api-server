# Build settings
BINARY_NAME=app
CMD_DIR=./cmd
BUILD_DIR=./bin
MAIN_FILE=$(CMD_DIR)/main.go
DOCS_FILE=$(CMD_DIR)/docs/docs.go

# Create Swagger Docs
swagger:
	@swag fmt
	@swag init -d cmd,internal/http --parseDependency

# Build the application
build:
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Clean build artifacts
clean:
	@rm -rf $(BUILD_DIR)/$(BINARY_NAME)

# Run the application
run:
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run the application
run-dev:
	@go run $(MAIN_FILE)

# Run the application
run-watch:
	@air --build.cmd "go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)" --build.bin "$(BUILD_DIR)/$(BINARY_NAME)" --build.exclude_dir ".angular,node_modules,src,swagger-ui"
