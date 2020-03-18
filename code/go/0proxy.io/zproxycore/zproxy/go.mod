module zproxy

require (
	0proxy.io/core v0.0.0
	0proxy.io/zproxycore v0.0.0
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/mattn/go-runewidth v0.0.8 // indirect
	github.com/spf13/viper v1.6.2
)

replace 0proxy.io/core => ../../core

replace 0proxy.io/zproxycore => ../../zproxycore

go 1.13
