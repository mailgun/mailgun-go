.PHONY: all
.DEFAULT_GOAL := all

gen:
	rm events/events_easyjson.go
	easyjson --all events/events.go
	rm events/objects_easyjson.go
	easyjson --all events/objects.go

all:
	export GO111MODULE=on; go test . -v
