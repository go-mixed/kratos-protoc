# named middleware for kratos v2.6.2

The main code forked from [kratos/protoc-gen-go-http](https://github.com/go-kratos/kratos/tree/main/cmd/protoc-gen-go-http), 
So you can build `xxx_http.pb.go` via this package

The version is always following `kratos/protoc-gen-go-http`

**NO NEED** to install `kratos/protoc-gen-go-http`.

# Usage

In your kratos project:

```golang
import (
    "github.com/go-kratos/kratos/v2/transport/http"
	"gopkg.in/go-mixed/protoc-gen-go-middleware/named"
)

svr := http.Server(
    http.Filter(named.EnableMiddleware()),
	http.Middleware(named.KratosMiddleware("name", middleware1)),
	...
)

```

# Install for protoc

```bash
go install gopkg.in/go-mixed/protoc-gen-go-middleware
```

# How to Build `xxx_http.pb.go`

1. Install [protoc](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation)
2. Install `protoc-gen-go`
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
3. Build & Install `protoc-gen-go-middleware`(this package)

4. Build your proto(see [examples/test.proto](examples/test.proto)

```bash
protoc --proto_path=./ --proto_path=/usr/include \ 
  --go_out=paths=source_relative:. \
  --go-middleware_out=paths=source_relative:. \  
  examples/test.proto
```
> **NO NEED** `--go-http_out=`

# Development

If you modified the `pb/middleware.proto`, then compile the proto file of `pb/middleware.proto`. Then reinstall this project

```bash
protoc --proto_path=./ \
  --proto_path=/usr/include \
  --go_out=paths=source_relative:. \
  pb/middleware.proto
```