REPO_NAME=github.com/echoturing.tools



.PHONY: fmt
fmt:
	@find . -name "*.go" | xargs goimports -w -l --local $(REPO_NAME) --private "mockprivate"


.PHONY: test
test:
	@go test


.PHONY: build
build:
	@go build ../.