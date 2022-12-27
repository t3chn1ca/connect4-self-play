module github/connect4-selfplay

replace api => ./src/api

replace db => ./src/db

replace mcts => ./src/mcts

replace proto => ./src/proto

replace shared => ./src/shared

go 1.16

require (
	api v0.0.0-00010101000000-000000000000
	db v0.0.0-00010101000000-000000000000
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/gin-gonic/gin v1.7.1
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/pprof v0.0.0-20221219190121-3cb0bae90811 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20220517205856-0058ec4f073c // indirect
	golang.org/x/sys v0.3.0 // indirect
	google.golang.org/grpc v1.37.0
	mcts v0.0.0-00010101000000-000000000000
	proto v0.0.0-00010101000000-000000000000
	shared v0.0.0-00010101000000-000000000000
)
