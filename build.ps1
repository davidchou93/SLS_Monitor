dep ensure -v
$env:GOOS="linux"
go build -ldflags="-s -w" -o bin/echo echo/main.go
go build -ldflags="-s -w" -o bin/receiver receiver/main.go
go build -ldflags="-s -w" -o bin/deliever deliever/main.go