module github.com/mailgun/mailgun-go/v3

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/drhodes/golorem v0.0.0-20160418191928-ecccc744c2d9
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-ini/ini v1.41.0 // indirect
	github.com/gobuffalo/envy v1.6.11
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/jtolds/gls v4.2.1+incompatible // indirect
	github.com/mailru/easyjson v0.0.0-20180823135443-60711f1a8329
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/pkg/errors v0.8.1
	github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d // indirect
	github.com/smartystreets/goconvey v0.0.0-20181108003508-044398e4856c // indirect
	github.com/thrawn01/args v0.3.0
	gopkg.in/ini.v1 v1.41.0 // indirect
)

replace github.com/mailgun/mailgun-go/v3/events => ./events

replace github.com/mailgun/mailgun-go/v3/cmd/mailgun => ./cmd/mailgun
