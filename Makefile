build:
	@echo 'Building pricing service...'
	go build -ldflags='-s' -o=./bin/crm ./cmd/api

run: build
	./bin/crm
