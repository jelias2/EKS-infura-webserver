module github.com/jelias2/infra-test

go 1.16

require (
	github.com/go-resty/resty/v2 v2.6.0
	github.com/gorilla/mux v1.8.0
	github.com/jelias2/infra-test/apis v0.0.0 // indirect
	github.com/jelias2/infra-test/handlers v0.0.0
	github.com/sirupsen/logrus v1.8.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.0
)

replace (
	github.com/jelias2/infra-test/apis => ./apis
	github.com/jelias2/infra-test/handlers => ./handlers
)
