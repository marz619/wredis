version=0.1.5

.PHONY: all

all:
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  build         - build the dist binary"
	@echo "  clean         - clean the dist build"
	@echo "  deps          - pull and setup dependencies"
	@echo "  test          - run tests"
	@echo "  update_deps   - update go mod

build: clean
	@go build ./...
	@go vet ./...
	@golint ./...

clean:
	@rm -rf ./bin

deps:
	@go mod download

cover:
	@ginkgo -r -cover
	@go tool cover -html=wredis.coverprofile

test:
	@ginkgo -r -v -cover -race

update_deps:
	@go mod verify
	@go mod tidy
