# protoc plugin for kratos v2

The main code forked from [kratos/protoc-gen-go-http](https://github.com/go-kratos/kratos/tree/main/cmd/protoc-gen-go-http), 
So you can build `xxx_http.pb.go` via this package

The version is always following `kratos/protoc-gen-go-http`

#  Usage

- [**Named middleware**](docs/named_middleware.md): 
  - Call the middleware by name in the `api/xxx.proto` file
  - Support multiple middleware
  - Support middleware arguments
- [**Custom Request/Response for http**](docs/custom.md): 
  - Custom request
    - Parse the request body with your own logic
  - Custom response
    - Write custom HEAD/BODY to http response 
    - Or write streaming content to http response (SSE/Chunked/Download)
  - Not for gRPC
- [**File upload**](docs/upload.md)
  - Support upload file as request body
  - Support upload file via `multipart/form-data` request

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
import "http/options.proto";

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

// stream response example
rpc Download(DownloadRequest) returns (EmptyResponse) {
    option (google.api.http) = {
      post: "/v1/download",
      body: "*",
    };
    option (http.options) = {
      custom_response: true,
    };
}

// upload request example
rpc Upload(EmptyResponse) returns (UploadResponse)
    option (google.api.http) = {
      post: "/v1/upload",
      body: "*",
    };
    option (http.options) = {
      custom_request: true,
    };
}

```

# Development

If you modified the `middleware/middleware.proto` or `http/options.proto`, 
you MUST recompile it.

```bash
cd protoc-gen-go-http

protoc --proto_path=./ \
  --proto_path=./protoc-gen-go-http/pb \
  --go_out=paths=source_relative:. \
  pb/middleware/middleware.proto
  
protoc --proto_path=./ \
  --proto_path=./protoc-gen-go-http/pb \
  --go_out=paths=source_relative:. \
  pb/http/options.proto
```

Then, manually install the `protoc-gen-go-http`

```bash
go install .
```