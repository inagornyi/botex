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

# Colors
GREEN_COLOR   = \033[0;32m
DEFAULT_COLOR = \033[m

help:
	@echo 'Usage: make <TARGET>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    clean              remove binary'
	@echo '    test               test packages'	
	@echo '    format             reformat package sources'
	@echo '    build              compile packages and dependencies'
	@echo '    run                compile and run Go program'
	@echo ''	

clean:
	@echo "$(GREEN_COLOR)[CLEAN]$(DEFAULT_COLOR)"
	@$(GOCLEAN)
	@if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi

test:
	@echo "$(GREEN_COLOR)[$(TEST_STRING)]$(DEFAULT_COLOR)"
	@$(GOTEST)

format:
	@echo "$(GREEN_COLOR)[FORMAT]$(DEFAULT_COLOR)"
	@$(GOFMT)

build:
	@echo "$(GREEN_COLOR)[BUILD]$(DEFAULT_COLOR)"
	@$(GOBUILD) $(LDFLAGS) -o $(BINARY)

run:
	@echo "$(GREEN_COLOR)[RUN]$(DEFAULT_COLOR)"
	@$(GORUN) main.go