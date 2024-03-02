# Upload Request

## Feature

- [Request Body as file raw data](#request-body-as-file-raw-data)
- [Upload a file with the `multipart/form-data` request](#upload-a-file-with-the-multipartform-data-request)

## Request Body as file raw data

All the file content is in the request body.

1. Copy proto files to your project

- `protoc-gen-go-http/pb/http/options.proto` -> `your_project/third_party/http/options.proto`

2. set "custom_request" to `true` in the `api/xxx.proto` file

```proto
syntax = "proto3";

import "google/api/annotations.proto";
import "http/options.proto";

service Test {
    rpc Upload(UploadRequest) returns (UploadResponse) {
        option (google.api.http) = {
          post: "/v1/upload",
          body: "*",
        };
        option (http.options) = {
          custom_request: true,
        };
    }
}

message UploadRequest {
    string file_name = 1;
}

message UploadResponse {
    string message = 1;
}

```

3. Generate `xxx_http.pb.go`

See: [Generate xxx_http.pb.go](../README.md#generate-xxx_http.pb.go)

4. Service code example

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


func (s *TestService) Download(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	// get the http context of kratos
	httpCtx := ctx.Value("httpContext").(kratosHttp.Context)
	
	file, err := os.Create("your_file_path/" + req.FileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	// 
	if err = io.Copy(httpCtx.Request().Body, file); err != nil {
		return nil, err
    }
	
    return &pb.UploadResponse{
		Message: "success",
    }, nil
}
```

5. Call the service

```http request
POST /v1/upload?file_name=the_file_name HTTP/1.1
Host: localhost:
Content-Type: application/octet-stream
Content-Length: 1234

<file raw data>
```


## Upload a file with the `multipart/form-data` request.

Coming soon...
