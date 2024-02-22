
# Named middleware

You can call the middleware by name in the `api/xxx.proto` file.
And define the named middleware in the `http.Server` of kratos

## 1. Copy proto files to your project

- `protoc-gen-go-http/pb/middleware/middleware.proto` -> `your_project/third_party/middleware/middleware.proto`


## 2. Write `xxx.proto` like this:

> See example： [test.proto](../protoc-gen-go-http/examples/test.proto)

Add the `middleware.caller` option to the rpc

```proto
syntax = "proto3";

import "google/api/annotations.proto";
import "middleware/middleware.proto";

// middleware.caller example
rpc User(UserRequest) returns (UserResponse) {
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

- the name of middleware is case-sensitive
- If you set the arguments, you must use "Type 3", It'll be ignored if you use other types

## 3. Generate `xxx_http.pb.go`

See: [Generate xxx_http.pb.go](../README.md#generate-xxx_http.pb.go)

## 4. Enable the named middleware in `http.Server`
**MUST** enable the named middleware before using it.

```golang
http.Filter(namedMiddlware.EnableMiddleware())
```

## 5. Register the named middleware in `http.Server`, 3 types of middleware are supported


### Type 1. Register a named middleware of kratos

```golang
http.Server(
    http.Filter(namedMiddlware.EnableMiddleware()), // enable the named middleware
    http.Middleware(namedMiddlware.KratosMiddleware("auth", func(nextHandler middleware.Handler) middleware.Handler{
       return func(ctx context.Context, req interface{}) (interface{}, error) {
           // do something
           return nextHandler(ctx, req)
       }
    }),
)
```

### Type 2. register a named handler of middleware of kratos

```golang
http.Server(
    http.Filter(namedMiddlware.EnableMiddleware()), // enable the named middleware
    http.Middleware(namedMiddlware.KratosHandler("auth", func(ctx context.Context, req interface{}) (interface{}, error) {
       return something, nil
    })),
)
   
```
- if error is not nil, the next handler will not be called, and the error will be returned to the client

### Type 3. register a named handler of middleware with arguments

```golang
http.Server(
    http.Filter(namedMiddlware.EnableMiddleware()), // enable the named middleware
    http.Middleware(namedMiddlware.HandlerWithArguments("auth", func(ctx context.Context, req interface{}, args ...string) (interface{}, error) {
       if len(args) > 0 {
           fmt.Println(args[0]) 
       }
       return something, nil
	})),
)
```

- if error is not nil, the next handler will not be called, and the error will be returned to the client
- the arguments in the `api/xxx.proto` will be passed to the middleware when the middleware is called


## 6. Example:

```golang 
   import (
       "context"
       "github.com/go-kratos/kratos/v2/transport/http"
       "github.com/go-mixed/kratos-protoc/namedMiddlware"
       "github.com/go-kratos/kratos/v2/middleware"
   )
   
   httpSrv := http.Server(
       http.Address(":8000")
       http.Filter(namedMiddlware.EnableMiddleware()),  // enable the named middleware
       http.Middleware(namedMiddlware.KratosMiddleware("auth", authMiddleware)), // register a named middleware of "auth"
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