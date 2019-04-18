export GO111MODULE=on
export GOPROXY=https://goproxy.io

PACKAGES=`go list ./... | grep -v /vendor/`
VETPACKAGES=`go list ./... | grep -v /vendor/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

default: deps 

list:
	@echo ${PACKAGES}
	@echo ${VETPACKAGES}
	@echo ${GOFILES}

lint:
	golangci-lint run

fmt:
	@gofmt -s -w ${GOFILES}

fmt-check:
	@diff=$$(gofmt -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

test:
	go test -cover -v ./...

vet:
	@go vet $(VETPACKAGES)

clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi

vendor:
	go mod vendor

get:
	go get -u ./...

tidy:
	go mod tidy

deps:
	go build -v ./...

upgrade:
	go get -u

.PHONY: default fmt fmt-check lint install test run vendor gen tidy vet clean
