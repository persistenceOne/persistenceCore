export GO111MODULE = on

build: go.sum
	go build -o ${GOBIN}/hubClient commands/hub/hubClient/main.go
	go build -o ${GOBIN}/hubNode commands/hub/hubNode/main.go

go.sum: go.mod
	@echo "--> Verify Dependency Modification"
	@go mod verify

.PHONY: all build test benchmark