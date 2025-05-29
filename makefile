PLUGIN_DIR := plugins

# All Go files under plugins/*/
ALL_GO_SRCS := $(wildcard $(PLUGIN_DIR)/*/*.go)

# Only keep files where the filename matches the folder name
PLUGIN_SRCS := $(shell \
	for file in $(ALL_GO_SRCS); do \
		dir=$$(basename $$(dirname $$file)); \
		base=$$(basename $$file .go); \
		if [ "$$dir" = "$$base" ]; then echo $$file; fi; \
	done)

# Substitute .go files for binaries
PLUGIN_BINS := $(patsubst %.go, %, $(PLUGIN_SRCS))

.PHONY: all
all: build-plugins

.PHONY: build-plugins
build-plugins: $(PLUGIN_BINS)

%:
	@echo "Building $@"
	@dir=$$(dirname $@); \
	go build -o $@ $$dir/*.go

.PHONY: clean
clean:
	@echo "Cleaning plugin binaries..."
	@for bin in $(PLUGIN_BINS); do \
		if [ -f $$bin ]; then rm -v $$bin; fi; \
	done

.PHONY: print-debug
print-debug:
	@echo "ALL_GO_SRCS = $(ALL_GO_SRCS)"
	@echo "PLUGIN_SRCS = $(PLUGIN_SRCS)"
	@echo "PLUGIN_BINS = $(PLUGIN_BINS)"
