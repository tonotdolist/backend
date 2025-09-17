cd ../
go test ./... -coverprofile=coverage ./...
go tool cover -html coverage
pause