.DEFAULT_GOAL := all

PACKAGE := github.com/mailgun/mailgun-go
GOPATH=$(shell go env GOPATH)
# TODO(vtopc): move into mtypes/internal/...?
TYPES_PATH=./internal/types

GOLANGCI_LINT_VERSION=v2.10.1
GOLANGCI_LINT_PATH=$(GOPATH)/bin/golangci-lint-v2
GOLANGCI_LINT=$(GOLANGCI_LINT_PATH)/golangci-lint

NILAWAY = $(GOPATH)/bin/nilaway
$(NILAWAY):
	go install go.uber.org/nilaway/cmd/nilaway@latest

.PHONY: all
all: test

.PHONY: test
test:
	go test ./... -race -count=1

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
	$(NILAWAY) -include-pkgs="$(PACKAGE)" -test=false -exclude-errors-in-files=mocks/ ./...

# linter:
$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GOLANGCI_LINT_PATH) $(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

## Download OpenAPI 3.1 spec files and generate models
.PHONY: get-and-gen-models
get-and-gen-models: get-openapi convert-openapi gen-models

.PHONY: get-openapi
get-openapi:
	# https://documentation.mailgun.com/docs/mailgun/api-reference/send/mailgun
	curl -o $(TYPES_PATH)/mailgun/mailgun.yaml https://documentation.mailgun.com/_spec/docs/mailgun/api-reference/send/mailgun.yaml?download
	# https://documentation.mailgun.com/docs/inboxready/api-reference/optimize/inboxready
	curl -o $(TYPES_PATH)/inboxready/inboxready.yaml https://documentation.mailgun.com/_spec/docs/inboxready/api-reference/optimize/inboxready.yaml?download

## Downgrade openapi 3.1 to 3.0
# this is one of the official ways to support OpenAPI 3.1:
# https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#does-oapi-codegen-support-openapi-31
# but it doesn't support `anyOf: [{type}, null]` for nullable fields -
# https://www.jvt.me/posts/2025/05/04/oapi-codegen-trick-openapi-3-1/
#
# install openapi-down-convert:
#  npm i -g @apiture/openapi-down-convert
#
# TODO(vtopc): use https://github.com/oapi-codegen/oapi-codegen-exp instead, which supports OpenAPI 3.1?
#	module declares its path as: github.com/oapi-codegen/oapi-codegen/v2
#		but was required as: github.com/oapi-codegen/oapi-codegen-exp/v2
# TODO(Go1.24): move into tools of go.mod(https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#for-go-124)?
# TODO(v6): switch to https://github.com/doordash-oss/oapi-codegen-dd instead?
# install oapi-codegen:
#  go install github.com/oapi-codegen/oapi-codegen-exp/v2/cmd/oapi-codegen@latest
#
# ValidateEmailResponse is described here better, than in the OpenAPI documentation, so we are not generating it.
# TODO(v6?): call gen-mailgun-models
.PHONY: gen-models
gen-models: gen-inboxready-models

## Generate Mailgun Send models
.PHONY: gen-mailgun-models
gen-mailgun-models:
	oapi-codegen -config $(TYPES_PATH)/mailgun/codegen_cfg.yaml $(TYPES_PATH)/mailgun/openapi_3.0.yaml

## Generate Mailgun Optimize models
.PHONY: gen-inboxready-models
gen-inboxready-models:
	oapi-codegen -config $(TYPES_PATH)/inboxready/codegen_cfg.yaml $(TYPES_PATH)/inboxready/openapi_3.0.yaml
