# func_gin_vmess_ping
#### go build

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build

#### linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
#### mac
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build  
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build

#### windows
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build 

#####
http://127.0.0.1:8080/vmess_ping?vmess=eyJhZGQiOiIxMy4yMjkuMTI3LjE4NCIsInBhdGgiOiIiLCJwcyI6ImF3cy1rdGNwIiwicG9ydCI6IjMzMDYiLCJ2IjoiMiIsImhvc3QiOiIiLCJ0bHMiOiIiLCJpZCI6IjVkNDg5M2EwLTE4ZDUtMTFlYi1hNTAxLTAyOTQwNWJiOTIwZSIsIm5ldCI6ImtjcCIsInR5cGUiOiJub25lIiwiYWlkIjoiMiIsInNuaSI6IiJ9

http://48454388-1356827337907157.test.functioncompute.com/vmess_ping?vmess=eyJhZGQiOiIxMy4yMjkuMTI3LjE4NCIsInBhdGgiOiIiLCJwcyI6ImF3cy1rdGNwIiwicG9ydCI6IjMzMDYiLCJ2IjoiMiIsImhvc3QiOiIiLCJ0bHMiOiIiLCJpZCI6IjVkNDg5M2EwLTE4ZDUtMTFlYi1hNTAxLTAyOTQwNWJiOTIwZSIsIm5ldCI6ImtjcCIsInR5cGUiOiJub25lIiwiYWlkIjoiMiIsInNuaSI6IiJ9



#### go mod download

#### go mod tidy
