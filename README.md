# protoc plugin for kratos v2

The main code forked from [kratos/protoc-gen-go-http](https://github.com/go-kratos/kratos/tree/main/cmd/protoc-gen-go-http), 
So you can build `xxx_http.pb.go` via this package

The version is always following `kratos/protoc-gen-go-http`


# Prerequisites 

1. Install [protoc](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation)
2. Install `protoc-gen-go`
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    ```
3. Build & Install `protoc-gen-go-http`
   ```bash
   go install github.com/go-mixed/kratos-protoc/protoc-gen-go-http@latest
   ```

   **NO NEED** to install official `kratos/protoc-gen-go-http`.

#  Usage

- [Named middleware](docs/named_middleware.md): Call the middleware by name in the `api/xxx.proto` file
- [Stream response](docs/stream.md): Write streaming content to http response (SSE/Chunked)

# Generate xxx_http.pb.go

It's no different from the official `kratos/protoc-gen-go-http`

```bash
kratos proto client ./your_project/api/v1/examples/test.proto
```

or

```bash
protoc --proto_path=./ \
   --proto_path=./ \
   --proto_path=./third_party \
   --go_out=paths=source_relative:. \
   --go-grpc_out=paths=source_relative:. \
   --go-http_out=paths=source_relative:. \
   --openapi_out=paths=source_relative:. \
   ./your_project/api/v1/examples/test.proto
 ```

# Proto file example
```proto
syntax = "proto3";

import "google/api/annotations.proto";
import "middleware/middleware.proto";
import "response/response.proto";

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

// stream.response example
rpc OpenAI(OpenAIRequest) returns (OpenAIResponse) {
    option (google.api.http) = {
      post: "/v1/openai",
      body: "*",
    };
    option (response.options) = {
      custom: true,
    };
}
```

# Development

If you want to modify the `protoc-gen-go-http/pb/middleware/middleware.proto` or `response/response.proto`, 
you MUST recompile it. 

```bash
protoc --proto_path=./ \
  --proto_path=./protoc-gen-go-http/pb \
  --go_out=paths=source_relative:. \
  protoc-gen-go-http/pb/middleware/middleware.proto
  
protoc --proto_path=./ \
  --proto_path=./protoc-gen-go-http/pb \
  --go_out=paths=source_relative:. \
  protoc-gen-go-http/pb/response/response.proto
```

Then, manually install the `protoc-gen-go-http`

```bash
cd protoc-gen-go-http 
go install .
```