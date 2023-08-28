BINARY_NAME=vuyo

build:
	GOARCH=amd64 GOOS=darwin go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o build/${BINARY_NAME}-darwin vuyo.go
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o build/${BINARY_NAME}-linux vuyo.go
	GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -gcflags "all=-trimpath=$(pwd)" -o build/${BINARY_NAME}-windows.exe vuyo.go

clean:
	rm -rf build/${BINARY_NAME}_*