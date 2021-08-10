module github.com/jelias2/infra-test

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	github.com/sirupsen/logrus v1.8.1
	github.com/jelias2/infra-test/handlers v0.0.0
	github.com/jelias2/infra-test/apis v0.0.0
)

replace (
	github.com/jelias2/infra-test/handlers => ./handlers
	github.com/jelias2/infra-test/apis => ./apis
)
