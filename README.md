# named middleware for kratos v2.6.2

The main code forked from [kratos/protoc-gen-go-http](https://github.com/go-kratos/kratos/tree/main/cmd/protoc-gen-go-http), 
So you can build `xxx_http.pb.go` via this package

The version is always following `kratos/protoc-gen-go-http`

**NO NEED** to install `kratos/protoc-gen-go-http`.

# Usage

## In your kratos project:

```golang
import (
    "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-mixed/kratos-middleware/named"
    "github.com/go-kratos/kratos/v2/middleware"
)

svr := http.Server(
    http.Filter(named.EnableMiddleware()),
	http.Middleware(named.KratosMiddleware("auth", authMiddleware)),
	...
)

func authMiddleware(next middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
        // do something
        return next(ctx, req)
    }
}

```

## proto file:

See [examples/test.proto](protoc-gen-go-http/examples/test.proto)

```proto
rpc User(UserRequest) returns (Response) {
    option (google.api.http) = {
      post: "/v1/user",
      body: "*",
    };
    // 这里的中间件顺序不是真正的执行顺序，真正的顺序在创建`http.Server`时候就已经确定了，比如：`http.Server(http.Middleware(middleware1, middleware2, middleware3))`
    // the order of middleware in this place is NOT guaranteed, the real order is determined when creating `http.Server`, for example: `http.Server(http.Middleware(middleware1, middleware2, middleware3))`
    option (pb.middleware) = {
      name: "auth",
      arguments: "bearer"
      arguments: "token"
    };
}
```


# Install for protoc module

```bash
go install github.com/go-mixed/kratos-middleware/protoc-gen-go-http
```

# How to Build `xxx_http.pb.go`

1. Install [protoc](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation)
2. Install `protoc-gen-go`
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    ```
3. Build & Install `protoc-gen-go-middleware`(see above)

4. Build your proto(see [examples/test.proto](protoc-gen-go-http/examples/test.proto)

    ```bash
    protoc --proto_path=./ \
      --proto_path=./protoc-gen-go-http/pb \
      --proto_path=/usr/include \
      --go_out=paths=source_relative:. \
      --go-http_out=paths=source_relative:. \
      protoc-gen-go-http/examples/test.proto
    ```

# Development

If you modified the `protoc-gen-go-middleware/pb/middleware.proto`, then compile the proto file of `pb/middleware.proto`. Then reinstall this project

```bash
protoc --proto_path=./ \
  --proto_path=./protoc-gen-go-http/pb \
  --proto_path=/usr/include \
  --go_out=paths=source_relative:. \
  protoc-gen-go-http/pb/middleware/middleware.proto
```