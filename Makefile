.DEFAULT_GOAL := all

PACKAGE := github.com/mailgun/mailgun-go
GOPATH=$(shell go env GOPATH)
TYPES_PATH=./internal/types

NILAWAY = $(GOPATH)/bin/nilaway
$(NILAWAY):
	go install go.uber.org/nilaway/cmd/nilaway@latest

.PHONY: all
all: test

.PHONY: test
test:
	go test . -race -count=1

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
	$(GOLINT) run

# go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
#
# mailgun/api-reference/openapi-final.yaml fails due to interface{} fields
#	# generate mailgun models
#	cd $(TYPES_PATH)/redocly-mailgun/docs/mailgun/api-reference/ && sed -i '' 's/openapi: 3.1.0/openapi: 3.0.0/' openapi-final.yaml
#	oapi-codegen -config $(TYPES_PATH)/mailgun_cfg.yaml $(TYPES_PATH)/redocly-mailgun/docs/mailgun/api-reference/openapi-final.yaml
.PHONY: gen-models
gen-models:
	cd $(TYPES_PATH) && git clone --depth 1 git@github.com:mailgun/redocly-mailgun.git
	# generate inboxready models
	sed -i '' 's/openapi: 3.1.0/openapi: 3.0.0/' $(TYPES_PATH)/redocly-mailgun/docs/inboxready/api-reference/openapi-final.yaml
	oapi-codegen -config $(TYPES_PATH)/inboxready_cfg.yaml $(TYPES_PATH)/redocly-mailgun/docs/inboxready/api-reference/openapi-final.yaml
	# generate validate models
	sed -i '' 's/openapi: 3.1.0/openapi: 3.0.0/' $(TYPES_PATH)/redocly-mailgun/docs/inboxready/api-reference/openapi-validate-final.yaml
	oapi-codegen -config $(TYPES_PATH)/validate_cfg.yaml $(TYPES_PATH)/redocly-mailgun/docs/inboxready/api-reference/openapi-validate-final.yaml
	rm -rf $(TYPES_PATH)/redocly-mailgun
