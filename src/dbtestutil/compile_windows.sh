#bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o cacheBench-windows.exe main.go