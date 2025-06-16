include .env

envs:
	@echo "  >  $(PACKAGE)"
	@echo "  >  $(VERSION)"
	@echo "  >  $(COMMIT_HASH)"
	@echo "  >  $(BUILD_TIMESTAMP)"
	@echo "  >  $(LDFLAGS)"


test:
	go test ./...

build:
	go build -ldflags "$(LDFLAGS)" -o bin/main main.go

deps:
	go get

run:
	go run main.go