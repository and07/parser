BUILD_ENVPARMS:=CGO_ENABLED=0
BIN?=./bin/parser
VGO_EXEC:=go
export GO111MODULE=on
APP?=parser

# install project dependencies
.PHONY: .deps
deps:
	$(info #Install dependencies...)
	$(VGO_EXEC) mod download

# run unit tests
.PHONY: .test
test: deps
	$(info #Running tests...)
	$(VGO_EXEC) test ./... -cover

.PHONY: .fast-build
fast-build: deps
	$(info #Building...)
	$(BUILD_ENVPARMS) $(VGO_EXEC) build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/${APP}

.PHONY: .build
build: test fast-build

run:
	$(BUILD_ENVPARMS) $(VGO_EXEC) run -ldflags "$(LDFLAGS)" ./cmd/${APP}
