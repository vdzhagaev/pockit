include local.env
export

run_local:
	go fmt ./... && \
	go vet ./... && \
	@echo "CONFIG_PATH is $(CONFIG_PATH)"
	go run ./cmd/url-shortener/main.go