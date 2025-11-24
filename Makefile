build:
	go build -ldflags "-X 'github.com/dontagr/pandora/internal/agent/config.bildDT=$(shell date +%FT%T)'" -o cmd/agent/agent cmd/agent/main.go