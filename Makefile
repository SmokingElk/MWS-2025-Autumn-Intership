.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: build
build:
	go build -o ./update-checker.exe ./cmd/main.go

.PHONY: mocks
mocks:
	mockgen -source=internal/domain/module/interfaces/module-client.go \
	-destination=internal/domain/module/mocks/mock-module-client.go
	mockgen -source=internal/domain/repo/interfaces/repo-hub-client.go \
	-destination=internal/domain/repo/mocks/mock-repo-hub-client.go

.PHONY: test
test:
	go test -v -count=1 ./...

.PHONY: test10
test10:
	go test -v -count=10 ./...

.PHONY: cover
cover:
	go test -coverprofile=cover.out -count=1 ./...
	go tool cover -html=cover.out
	DEL cover.out

.PHONY: lint
lint:
	golangci-lint -v run ./...