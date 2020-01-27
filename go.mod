module github.com/mailgun/mailgun-go/v3

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/go-chi/chi v4.0.0+incompatible
	github.com/mailru/easyjson v0.7.0
	github.com/pkg/errors v0.8.1
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7 // indirect
	github.com/ugorji/go v1.1.7 // indirect
)

replace github.com/mailgun/mailgun-go/v3/events => ./events

go 1.13
