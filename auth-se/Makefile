#!/bin/bash

build: binary

binary:
	@echo "building binary.."
	@go build -ldflags "-X auth-se/cmd.buildVersion=`git rev-parse --short HEAD`" -tags static_all .


clean:
	@echo "cleaning ..."
	@rm -f auth-se
	@rm -rf vendor
	@rm -f go.sum


install:
	@echo "Installing dependencies...."
	@rm -rf vendor
	@rm -f Gopkg.lock
	@rm -f glide.lock
	@go mod tidy && go mod download && go mod vendor

test:
	@go test $$(go list ./... | grep -v /vendor/) -cover

test-cover:
	@go test $$(go list ./... | grep -v /vendor/) -coverprofile=cover.out && go tool cover -html=cover.out ; rm -f cover.out

coverage:
	@go test -covermode=count -coverprofile=count.out fmt; rm -f count.out

start:
	@go run -ldflags "-X auth-se/cmd.buildVersion=`git rev-parse --short HEAD`" main.go http

format:
	@go fmt $$(go list ./... | grep -v /vendor/)