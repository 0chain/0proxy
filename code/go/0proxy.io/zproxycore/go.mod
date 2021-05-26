module 0proxy.io/zproxycore

require (
	0proxy.io/core v0.0.0
	github.com/0chain/gosdk v1.2.5
	github.com/fatih/color v1.9.0 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-runewidth v0.0.8 // indirect
	go.uber.org/zap v1.15.0
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

replace 0proxy.io/core => ../core

go 1.13
