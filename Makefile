# Basic go commands
GOCMD     = go
GOBUILD   = $(GOCMD) build
GORUN     = $(GOCMD) run
GOCLEAN   = $(GOCMD) clean
GOTEST    = $(GOCMD) test
GOGET     = $(GOCMD) get
GOFMT     = $(GOCMD) fmt
GOVET     = $(GOCMD) tool vet

# Binary output name
BINARY = botex

# Basic docker commands
DOCKER       = docker
DOCKERBUILD  = $(DOCKER) build
DOCKERRUN    = $(DOCKER) run
DOCKERSTOP   = $(DOCKER) stop
DOCKERREMOVE = $(DOCKER) rm
DOCKERPS     = $(DOCKER) ps

# These will be provided to the target
VERSION = 1.0.0
BUILD   = $(shell git rev-parse HEAD)

# Setup the -ldflags option for go build here, interpolate the variable vulues
LDFLAGS = -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

#
BUILD_DIR    = $(CURDIR)
COVERAGE_DIR = $(BUILD_DIR)/coverage

#
PKGS = $(shell go list ./... | grep -v /vendor)

# Colors
GREEN_COLOR   = \033[0;32m
DEFAULT_COLOR = \033[m

# Texts
TEST_STRING = "TEST"

.PHONY: all help clean test coverage vet format build run docker version

all: clean format vet coverage test

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    clean              remove binary'
	@echo '    test               test packages'
	@echo '    coverage           report code tests coverage'
	@echo '    vet                report likely mistakes in packages'
	@echo '    format             reformat package sources'
	@echo '    build              compile packages and dependencies'
	@echo '    run                compile and run Go program'
	@echo ''
	@echo 'Targets run by default are: clean format vet coverage test'

clean:
	@echo "$(GREEN_COLOR)[CLEAN]$(DEFAULT_COLOR)"
	@$(GOCLEAN)
	@if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi

test:
	@echo "$(GREEN_COLOR)[$(TEST_STRING)]$(DEFAULT_COLOR)"
	@$(GOTEST) -v $(PKGS)

coverage:
	@echo "$(GREEN_COLOR)[COVERAGE]$(DEFAULT_COLOR)"
	@# Create the coverage files directory
	@mkdir -p $(COVERAGE_DIR)
	@# Create a coverage file for each package
	@for package in $(PKGS); do $(GOTEST) -covermode=count -coverprofile $(COVERAGE_DIR)/`basename "$$package"`.cov "$$package"; done
	@# Merge the coverage profile files
	@echo 'mode: count' > $(COVERAGE_DIR)/coverage.cov ;
	@tail -q -n +2 $(COVERAGE_DIR)/*.cov >> $(COVERAGE_DIR)/coverage.cov ;
	@# Display the global code coverage
	@go tool cover -func=$(COVERAGE_DIR)/coverage.cov ;
	@# If needed, generate HTML report
	@if [ $(html) ]; then go tool cover -html=$(COVERAGE_DIR)/coverage.cov -o coverage.html ; fi
	@# Remove the coverage files directory
	@rm -rf $(COVERAGE_DIR);

vet:
	@echo "$(GREEN_COLOR)[VET]$(DEFAULT_COLOR)"
	@-$(GOVET) -v $(PKGS)

format:
	@echo "$(GREEN_COLOR)[FORMAT]$(DEFAULT_COLOR)"
	@$(GOFMT) $(PKGS)

build:
	@echo "$(GREEN_COLOR)[BUILD]$(DEFAULT_COLOR)"
	@$(GOBUILD) $(LDFLAGS) -o $(BINARY)

run:
	@echo "$(GREEN_COLOR)[RUN]$(DEFAULT_COLOR)"
	@$(GORUN) main.go