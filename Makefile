NAME=scanner
VERSION=0.0.1

.PHONY: build
build:
	@go build -o $(NAME)

.PHONY: run
run: ensure build
	@./$(NAME) -e development

.PHONY: run-prod
run-prod: ensure build
	@./$(NAME) -e prod

.PHONY: clean
clean:
	@rm -f $(NAME)

.PHONY: ensure
ensure:
	@dep ensure --vendor-only

.PHONY: test
test: ensure
	@go test -v ./tests/*