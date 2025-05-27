PLUGIN_DIR := plugins

# Find all plugin source files matching ./plugins/<name>/<name>.go
PLUGIN_SRCS := $(wildcard $(PLUGIN_DIR)/*/*.go)

# Extract plugin names
PLUGIN_NAMES := $(basename $(notdir $(PLUGIN_SRCS)))

# Substitute .go files for binaries
PLUGIN_BINS := $(patsubst %.go, %, $(PLUGIN_SRCS))

.PHONY: all
all: build-plugins

.PHONY: build-plugins
build-plugins: $(PLUGIN_BINS)

%: %.go
	@echo "Building $@"
	@go build -o $@ $<

.PHONY: clean
clean:
	@echo "Cleaning plugin binaries..."
	@for bin in $(PLUGIN_BINS); do \
		if [ -f $$bin ]; then rm -v $$bin; fi; \
	done
