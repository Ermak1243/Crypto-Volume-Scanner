start-app:
	go run cmd/app/main.go
start-docs-server:
	godoc -http=:6060
parallel: start-docs-server start-app