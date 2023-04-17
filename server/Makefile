.PHONY: test server

test:
	@find . -type f -name '*_test.go' -exec dirname {} \; | uniq | while read -r dir; do \
		echo "Running tests in $$dir"; \
		go test -short $$dir; \
	done

server:
	go run cmd/app/main.go
