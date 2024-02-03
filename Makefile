CFLAGS="-g -O2 -Wno-return-local-addr"
CC=x86_64-w64-mingw32-gcc
CXX=x86_64-w64-mingw32-g++

test:
	cd tests && \
	CGO_CFLAGS=${CFLAGS} \
	go test

run:
	CGO_CFLAGS=${CFLAGS} \
	go run main.go --init-db

build:
	CGO_CFLAGS=${CFLAGS} \
	go build -o a.out main.go

build_release:
	CGO_CFLAGS=${CFLAGS} \
	GIN_MODE=release \
	go build -o a.out main.go

build_win:
	CGO_CFLAGS=${CFLAGS} \
	GOOS=windows \
	GOARCH=amd64 \
	CGO_ENABLED=1 \
	CC=${CC} \
	CXX=${CXX} \
	go build -o a.out.exe
