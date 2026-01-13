include local.env
export

run_local:
	go fmt ./... && \
	go vet ./... && \
	go run ./cmd/url-shortener/main.go