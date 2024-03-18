
# Named middleware

You can call the middleware by name in the `api/xxx.proto` file.
And define the named middleware in the `http.Server` of kratos

## 1. Copy proto files to your project

`protoc-gen-go-http/pb/middleware/middleware.proto` -> `your_project/third_party/middleware/middleware.proto`

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

## 3. Generate `xxx_http.pb.go`

See: [Generate xxx_http.pb.go](../README.md#generate-xxx_http.pb.go)

## 4. Enable the named middleware in `http.Server`
**MUST** enable the named middleware before using it.

```golang
http.Filter(namedMiddleware.EnableMiddleware())
```

## 5. Register the named middleware in `http.Server`, 2 methods of middleware are supported


### Type 1. Wrap a kratos middleware to a named middleware

```golang
http.Server(
    http.Filter(namedMiddleware.EnableMiddleware()), // enable the named middleware
    http.Middleware(
        namedMiddleware.WrapKratosMiddleware("auth", func(nextHandler middleware.Handler) middleware.Handler{
           return func(ctx context.Context, req interface{}) (interface{}, error) {
               // do something
               return nextHandler(ctx, req)
           }
        }),
    )
)
```

### Type 2. register a named middleware

```golang
http.Server(
    http.Filter(namedMiddleware.EnableMiddleware()), // enable the named middleware
    http.Middleware(
        namedMiddleware.Middleware("auth", func(ctx context.Context, lastReq interface{}) (req interface{}, error) {
            return lastReq, nil
        })
    ),
)
   
```
- if error is not nil, the next handler will not be called, and the error will be returned to the client

```
type Handler func(ctx context.Context, lastReq interface{}) (req interface{}, err error)
```

- `lastReq` is previous request from the previous middleware
- `req` is the request to the next middleware, you may modify the request and return, or return "lastReq" directly

## 6. Get the arguments

```golang
http.Server(
    http.Filter(namedMiddleware.EnableMiddleware()), // enable the named middleware
    http.Middleware(
        namedMiddleware.Middleware("auth", func(ctx context.Context, lastReq interface{}) (req interface{}, err error) {
            arguments := namedMiddleware.GetArguments(ctx) // get the arguments

            return req, nil
        })
    ),
)
```

- `namedMiddleware.GetArguments(ctx)` will return the arguments of this named middleware who is calling
- **You can only get the arguments in the named middleware body**


## 7. Example:

```golang 
   import (
       "context"
       "github.com/go-kratos/kratos/v2/transport/http"
       "github.com/go-mixed/kratos-protoc/namedMiddleware"
       "github.com/go-kratos/kratos/v2/middleware"
   )
   
   httpSrv := http.Server(
       http.Address(":8000")
       http.Filter(namedMiddleware.EnableMiddleware()),  // enable the named middleware
       http.Middleware(
            namedMiddleware.WrapKratosMiddleware("auth", authMiddleware) // register a named middleware
       ),
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
           arguments := namedMiddleware.GetArguments(ctx) // get the arguments
           // do something
           return next(ctx, req)
       }
   }
```