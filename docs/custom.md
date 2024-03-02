# Custom response/response

## Introduction
- Parse the request body with your own logic
- Write custom HEAD/BODY to http response

## Feature
if `"custom_request"` is true, ONLY parse the **Query/Url params** to the proto request,
and you can parse the request body with your own logic

if `"custom_response"` is true, you can write custom HEAD/BODY to http response,
and you **MUST** `return nil, nil`

## Usage

### 1. Copy proto files to your project

- `protoc-gen-go-http/pb/http/options.proto` -> `your_project/third_party/http/options.proto`

### 2. Write `xxx.proto` like this:

> See example： [test.proto](../protoc-gen-go-http/examples/test.proto)

Add the `http.options` option to the rpc

```proto
syntax = "proto3";

import "google/api/annotations.proto";
import "http/options.proto";

service Test {
    rpc Download(Empty) returns (SomeResponse) {
        option (google.api.http) = {
          post: "/v1/download",
          body: "*",
        };
        option (http.options) = {
          custom_request: true,
          custom_response: true,
        };
    }
}

message Empty {}

message SomeResponse {
    string data = 1;
}

```

### 3. Generate `xxx_http.pb.go`

See: [Generate xxx_http.pb.go](../README.md#generate-xxx_http.pb.go)


### 4. Service code example

Download file example

```golang
package service

import (
	"context"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"encoding/json"
	"os"
	"io"
)

type TestService struct {
}

func NewTestService() *TestService {
	return &TestService{}
}

type SomeRequest struct {
	FilePath string `json:"file_path"`
}

func (s *TestService) Download(ctx context.Context, _ *pb.Empty) (*pb.SomeResponse, error) {
	// get the http context of kratos
	httpCtx := ctx.Value("httpContext").(kratosHttp.Context)

	// if "custom_request" is true, you can get the request body,
	// and unmarshal it to a struct with your own logic
	requestBody, err := io.ReadAll(httpCtx.Request().Body)
	var req SomeRequest
	if err := json.Unmarshal(requestBody, &request); err != nil {
		return nil, err
	}

	file, err := os.Open(req.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	// if "custom_response" is true, you can write custom HEAD/BODY to http response
	// eg: stream the file to http response
	httpCtx.Response().Header().Set("Content-Disposition", "attachment; filename=example.mp4")
    if err = httpCtx.Stream(200, "video/mp4", file); err != nil {
        return nil, err
    } else {
		return nil, nil // return nil if you want to use custom response
	}
	
	// return json normally if you don't want to use Streaming response
	return &pb.SomeResponse{
		Data: ...,
	}, nil
}

```

- `return nil, nil`: the first "nil" means you want to use streaming response
- `return nil, err`: if err is not nil, the normal error response will be returned
