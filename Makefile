.PHONY: all
.DEFAULT_GOAL := all

gen:
	rm events/events_easyjson.go
	easyjson --all events/events.go
	rm events/objects_easyjson.go
	easyjson --all events/objects.go

all:
# Only run these tests if secure credentials exist
ifeq ($(TRAVIS_SECURE_ENV_VARS),true)
	go get -t .
	go test .
endif
