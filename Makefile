BINARY=internet-monitor

.PHONY: run build clean fmt test test_verbose
default: run

run: | clean
	@echo --------------------------------------------------
	@echo Running
	@echo --------------------------------------------------
	@go run main.go

build: | clean
	@echo --------------------------------------------------
	@echo Building
	@echo --------------------------------------------------
	@go build -o ${BINARY}

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

fmt:
	@echo --------------------------------------------------
	@echo Formatting
	@echo --------------------------------------------------
	@goimports -w .
	@gofmt -s -w .

test:
	@echo --------------------------------------------------
	@echo Testing
	@echo --------------------------------------------------
	@go test ./...

test_verbose:
	@echo --------------------------------------------------
	@echo Testing with Verbose Logging
	@echo --------------------------------------------------
	@go test -v ./...