build:
	go build -ldflags "-X 'github.com/dontagr/pandora/internal/agent/config.bildDT=$(shell date +%FT%T)'" -o cmd/agent/agent cmd/agent/main.go

build_win_32:
	GOOS=windows GOARCH=386 go build -ldflags "-X 'github.com/dontagr/pandora/internal/agent/config.bildDT=$(shell date +%FT%T)'" -o cmd/agent/agent.exe cmd/agent/main.go

build_win_64:
	GOOS=windows GOARCH=amd64 go build -ldflags "-X 'github.com/dontagr/pandora/internal/agent/config.bildDT=$(shell date +%FT%T)'" -o cmd/agent/agent64.exe cmd/agent/main.go

build_linux_arm:
	GOOS=linux GOARCH=arm go build -ldflags "-X 'github.com/dontagr/pandora/internal/agent/config.bildDT=$(shell date +%FT%T)'" -o cmd/agent/agent_linux_arm cmd/agent/main.go

build_linux_386:
	GOOS=linux GOARCH=386 go build -ldflags "-X 'github.com/dontagr/pandora/internal/agent/config.bildDT=$(shell date +%FT%T)'" -o cmd/agent/agent_linux_32 cmd/agent/main.go

build_linux_64:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X 'github.com/dontagr/pandora/internal/agent/config.bildDT=$(shell date +%FT%T)'" -o cmd/agent/agent_linux_64 cmd/agent/main.go

build_macos_386:
	GOOS=darwin GOARCH=386 go build -ldflags "-X 'github.com/dontagr/pandora/internal/agent/config.bildDT=$(shell date +%FT%T)'" -o cmd/agent/agent_macos_32 cmd/agent/main.go

build_macos_64:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'github.com/dontagr/pandora/internal/agent/config.bildDT=$(shell date +%FT%T)'" -o cmd/agent/agent_macos_64 cmd/agent/main.go

run_server:
	go build -o cmd/server/server cmd/server/main.go && cmd/server/server

build_server_win_32:
	GOOS=windows GOARCH=386 go build -o cmd/server/server cmd/server/main.go

build_server_win_64:
	GOOS=windows GOARCH=amd64 go build -o cmd/server/server cmd/server/main.go

build_server_linux_arm:
	GOOS=linux GOARCH=arm go build -o cmd/server/server cmd/server/main.go

build_server_linux_386:
	GOOS=linux GOARCH=386 go build -o cmd/server/server cmd/server/main.go

build_server_linux_64:
	GOOS=linux GOARCH=amd64 go build -o cmd/server/server cmd/server/main.go

build_server_macos_386:
	GOOS=darwin GOARCH=386 go build -o cmd/server/server cmd/server/main.go

build_server_macos_64:
	GOOS=darwin GOARCH=amd64 go build -o cmd/server/server cmd/server/main.go
