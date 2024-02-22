# Stream response

## 1. Copy proto files to your project

- `protoc-gen-go-http/pb/stream/stream.proto` -> `your_project/third_party/stream/stream.proto`

## 2. Write `xxx.proto` like this:

> See example： [test.proto](../protoc-gen-go-http/examples/test.proto)

Add the `stream.response` option to the rpc

```proto
syntax = "proto3";

import "google/api/annotations.proto";
import "stream/stream.proto";

// middleware.caller example
rpc User(UserRequest) returns (UserResponse) {
    option (google.api.http) = {
      post: "/v1/user",
      body: "*",
    };
    option (stream.response) = {
      enabled: true,
    };
}

```

## 3. Generate `xxx_http.pb.go`

See: [Generate xxx_http.pb.go](../README.md#generate-xxx_http.pb.go)


## 4. Service code example

```golang
package service

import (
	"context"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"net/http"
)

type TestService struct {
}

func NewTestService() *TestService {
	return &TestService{}
}

func (s *TestService) OpenAI(ctx context.Context, req *pb.OpenAIRequest) (*pb.OpenAIResponse, error) {
	// write streaming content to http response
	httpCtx := ctx.(kratosHttp.Context)

	response, err := http.Post(ctx, "https://api.openai.com/v1/chat/completions", "application/json", ...)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	
	// return nil if you want to use Streaming response
	if req.Stream {
		if err = httpCtx.Stream(200, "text/event-stream", response.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
	
	// return json normally
	return &pb.OpenAIResponse{
		Data: ...,
	}, nil
}

```

- `httpCtx.Stream(200, "text/event-stream", response.Body)`: write SSE content to http response
- `return nil, nil`: the first "nil" means you want to use streaming response
- `return nil, err`: if err is not nil, the normal error response will be returned
- `return &pb.OpenAIResponse{...}, nil`: the json of `pb.OpenAIResponse` will be returned