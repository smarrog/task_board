module github.com/smarrog/task-board/notification-service

go 1.25.4

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.12.0
	github.com/rs/zerolog v1.34.0
	github.com/smarrog/task-board/shared v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

replace github.com/smarrog/task-board/shared => ../shared
