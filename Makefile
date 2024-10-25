.DEFAULT_GOAL := all

PACKAGE := github.com/mailgun/mailgun-go
GOPATH=$(shell go env GOPATH)

NILAWAY = $(GOPATH)/bin/nilaway
$(NILAWAY):
	go install go.uber.org/nilaway/cmd/nilaway@latest

.PHONY: all
all: test

.PHONY: test
test:
	export GO111MODULE=on; go test . -v

.PHONY: godoc
godoc:
	mkdir -p /tmp/tmpgoroot/doc
	-rm -rf /tmp/tmpgopath/src/${PACKAGE}
	mkdir -p /tmp/tmpgopath/src/${PACKAGE}
	tar -c --exclude='.git' --exclude='tmp' . | tar -x -C /tmp/tmpgopath/src/${PACKAGE}
	echo -e "open http://localhost:6060/pkg/${PACKAGE}\n"
	GOROOT=/tmp/tmpgoroot/ GOPATH=/tmp/tmpgopath/ godoc -http=localhost:6060

# TODO(vtopc): fix mocks and enable nilaway for them too?
.PHONY: nilaway
nilaway: $(NILAWAY)
	$(NILAWAY) -include-pkgs="$(PACKAGE)" -test=false -exclude-errors-in-files=mock_ ./...

# linter:
GOLINT = $(GOPATH)/bin/golangci-lint
$(GOLINT):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.61.0

.PHONY: lint
lint: $(GOLINT)
	$(GOLINT) run -new-from-rev=master
