```text
$env:CGO_ENABLED="0" 
$env:GOOS="linux" 
$env:GOARCH="amd64" 
go build -v -o bin/vectorsql-server src/cmd/server.go
```

```
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -v -o bin/vectorsql-server src/cmd/server.go
```














