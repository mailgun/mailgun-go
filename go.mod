module github.com/mailgun/mailgun-go/v3

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/go-chi/chi v4.0.0+incompatible
	github.com/mailru/easyjson v0.0.0-20180823135443-60711f1a8329
	github.com/pkg/errors v0.8.1
)

replace github.com/mailgun/mailgun-go/v3/events => ./events

go 1.13
