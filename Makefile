build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/echo echo/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/receiver receiver/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/deliever deliever/main.go

.PHONY: deploy
deploy: build
	serverless deploy --verbose
