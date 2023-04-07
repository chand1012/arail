set dotenv-load

default:
  just --list --unsorted

run *params:
  go run main.go {{params}}

build:
  mkdir -p bin
  go build -v -o bin/arail main.go

test:
  go test ./...

clean:
  rm -rf bin
  go clean -cache
