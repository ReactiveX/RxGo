test:
	go clean -testcache ./...
	go test -race -timeout 15s ./... --tags=all
	go test -timeout 15s -run TestLeak
