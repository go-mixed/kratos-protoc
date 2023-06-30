# named middleware for kratos v2.6.2

The main code forked from [kratos/protoc-gen-go-http](https://github.com/go-kratos/kratos/tree/main/cmd/protoc-gen-go-http), 
So you can build `xxx_http.pb.go` via this package

The version is always following `kratos/protoc-gen-go-http`

**NO NEED** to install `kratos/protoc-gen-go-http`.

# Usage

1. Copy `protoc-gen-go-http/pb/middleware/middleware.proto` 

   to `your_project/third_party/middleware/middleware.proto`

2. enable and add the middleware to `Kratos Boot` :

   ```golang 
   import (
       "context"
       "github.com/go-kratos/kratos/v2/transport/http"
       "github.com/go-mixed/kratos-middleware/named"
       "github.com/go-kratos/kratos/v2/middleware"
   )
   
   httpSrv := http.Server(
       http.Address(":8000")
       http.Filter(named.EnableMiddleware()),  // enable the named middleware
       http.Middleware(named.KratosMiddleware("auth", authMiddleware)), // register a named middleware of "auth"
   )
   
   grpcSrv := grpc.NewServer(grpc.Address(":9000"))
   
   app := kratos.New(
      kratos.Name("kratos"),
      kratos.Version("latest"),
      kratos.Server(httpSrv, grpcSrv),
   )
   app.Run()
   
   func authMiddleware(next middleware.Handler) middleware.Handler {
       return func(ctx context.Context, req interface{}) (interface{}, error) {
           // do something
           return next(ctx, req)
       }
   }
   
   ```

3. build `api/xxx.proto` like this:

   See [examples/test.proto](protoc-gen-go-http/examples/test.proto)

   ```proto
   syntax = "proto3";
   
   import "google/api/annotations.proto";
   import "middleware/middleware.proto";
   
   rpc User(UserRequest) returns (Response) {
       option (google.api.http) = {
         post: "/v1/user",
         body: "*",
       };
       option (middleware.caller) = {
         name: "auth",
         arguments: "arg1"
         arguments: "arg2"
       };
   }
   ```

# API

add to the initialization of `http.Server`, Example:

   ```golang
   http.Server(
	   http.Filter(named.EnableMiddleware()), // enable the named middleware
       http.Middleware(named.KratosMiddleware("auth", xxxKratosMiddleware)), // register a middleware named "auth"
       http.Middleware(named.KratosHandler("auth", xxxKratosHandler)), // register a handler named "auth"
       http.Middleware(named.HandlerWithArguments("auth", xxxHandlerWithArguments)), // register a handler with arguments named "auth"
   )
   ```

## Enable the named middleware
MUST enable the named middleware before using it

   ```golang
   http.Filter(named.EnableMiddleware())
   ```
## Register middleware

### 1. Register a named middleware

   ```golang
   http.Middleware(named.KratosMiddleware("auth", func(nextHandler middleware.Handler) middleware.Handler{
       return func(ctx context.Context, req interface{}) (interface{}, error) {
           // do something
           return nextHandler(ctx, req)
       }
   })
   ```

### 3. Or a named handler of middleware

   ```golang
   http.Middleware(named.KratosHandler("auth", func(ctx context.Context, req interface{}) (interface{}, error) {
       return something, nil
   }))
   
   ```
- if error is not nil, the next handler will not be called, and the error will be returned to the client

### 4. Register a named handler of middleware with arguments

   ```golang
   http.Middleware(named.HandlerWithArguments("auth", func(ctx context.Context, req interface{}, args ...string) (interface{}, error) {
	   if len(args) > 0 {
           fmt.Println(args[0]) 
       }
       return something, nil
   }))
   ```
- if error is not nil, the next handler will not be called, and the error will be returned to the client
- the arguments in the `api/xxx.proto` will be passed to the middleware when the middleware is called 

## Call the named middleware

```protobuf
rpc ApiName(Request) returns (Response) {
  option (middleware.caller) = {
    name: "auth",
    arguments: "arg0"
    arguments: "arg1"
    };
}
```
- the name of middleware is case-sensitive
- the arguments will be ignored if the middleware is not registered with arguments

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

4. Build your proto(see [examples/test.proto](protoc-gen-go-http/examples/test.proto))

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