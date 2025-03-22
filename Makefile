.DEFAULT_GOAL := all

PACKAGE := github.com/mailgun/mailgun-go
GOPATH=$(shell go env GOPATH)
# TODO(vtopc): move into mtypes/internal/...?
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
	$(NILAWAY) -include-pkgs="$(PACKAGE)" -test=false -exclude-errors-in-files=mocks/ ./...

# linter:
GOLINT = $(GOPATH)/bin/golangci-lint
$(GOLINT):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.61.0

.PHONY: lint
lint: $(GOLINT)
	$(GOLINT) run

## Download OpenAPI 3.1 spec files and generate models
.PHONY: get-and-gen-models
get-and-gen-models: get-openapi convert-openapi gen-models

.PHONY: get-openapi
get-openapi:
	cd $(TYPES_PATH) && git clone --depth 1 git@github.com:mailgun/redocly-mailgun.git

## Convert openapi 3.1 to 3.0
# install openapi-down-convert:
#  npm i -g @apiture/openapi-down-convert
.PHONY: convert-openapi
convert-openapi:
	# Mailgun Send
	openapi-down-convert --input $(TYPES_PATH)/redocly-mailgun/docs/mailgun/api-reference/openapi-final.yaml --output $(TYPES_PATH)/mailgun/openapi_3.0.yaml
	# Mailgun Optimize
	openapi-down-convert --input $(TYPES_PATH)/redocly-mailgun/docs/inboxready/api-reference/openapi-final.yaml --output $(TYPES_PATH)/inboxready/openapi_3.0.yaml

# TODO(Go1.24): move into tools of go.mod(https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#for-go-124)?
# install oapi-codegen:
#  go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.2-0.20241128130830-b07f7ea6d520
#
# ValidateEmailResponse is described here better, than in the OpenAPI documentation, so we are not generating it.
# TODO(v6?): call gen-mailgun-models
.PHONY: gen-models
gen-models: gen-inboxready-models

## Generate Mailgun Send models
.PHONY: gen-mailgun-models
gen-mailgun-models:
	oapi-codegen -config $(TYPES_PATH)/mailgun/codegen_cfg.yaml $(TYPES_PATH)/mailgun/openapi_3.0.yaml
	# patch slices(`*[]` -> `[]`)
	sed -i '' 's/\*\[\]/\[\]/' $(TYPES_PATH)/mailgun/model.gen.go
	# patch maps(`*map` -> `map`)
	sed -i '' 's/\*map/map/' $(TYPES_PATH)/mailgun/model.gen.go

## Generate Mailgun Optimize models
.PHONY: gen-inboxready-models
gen-inboxready-models:
	oapi-codegen -config $(TYPES_PATH)/inboxready/codegen_cfg.yaml $(TYPES_PATH)/inboxready/openapi_3.0.yaml
	# TODO(vtopc): fix pointers to slices and maps with oapi-codegen's `x-go-type-skip-optional-pointer: true`?
	# patch slices(`*[]` -> `[]`)
	sed -i '' 's/\*\[\]/\[\]/' $(TYPES_PATH)/inboxready/model.gen.go
	# patch maps(`*map` -> `map`)
	sed -i '' 's/\*map/map/' $(TYPES_PATH)/inboxready/model.gen.go
