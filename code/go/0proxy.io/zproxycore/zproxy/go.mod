module zproxy

require (
	0proxy.io/core v0.0.0
	0proxy.io/zproxycore v0.0.0
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/spf13/viper v1.6.2
	go.uber.org/zap v1.15.0
)

replace 0proxy.io/core => ../../core

replace 0proxy.io/zproxycore => ../../zproxycore

go 1.13
