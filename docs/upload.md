# Upload Request

## Feature

- Request Body as file raw data
- Upload a file with the `multipart/form-data` request

## Example

All the file content is in the request body.

1. Copy proto files to your project

- `protoc-gen-go-http/pb/http/options.proto` -> `your_project/third_party/http/options.proto`

2. Set "custom_request" to `true` in the `api/xxx.proto` file

```proto
syntax = "proto3";

import "google/api/annotations.proto";
import "http/options.proto";

service Test {
    rpc Upload(Empty) returns (UploadResponse) {
        option (google.api.http) = {
          post: "/v1/upload",
          body: "*",
        };
        option (http.options) = {
          custom_request: true,
        };
    }
}

message Empty {
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


func (s *TestService) Upload(ctx context.Context, _ *pb.Empty) (*pb.UploadResponse, error) {
    // get the http context of kratos
    httpCtx := ctx.Value("httpContext").(kratosHttp.Context)

    // get the file name from the request
    file, err := upload.GetFileFromRequest(httpCtx.Request(), "file", 1024*1024*32)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    writer, err := os.Create("your_file_path/" + file.Name())
    if err != nil {
        return nil, err
    }
    defer writer.Close()
    
    // copy the file to the writer
    if err = io.Copy(file, writer); err != nil {
        return nil, err
    }
    
    return &pb.UploadResponse{
        Message: "success",
    }, nil
}
```

5. Call the service

```http request
POST /v1/upload?file=the_file_name HTTP/1.1
Host: localhost:
Content-Type: application/octet-stream
Content-Length: 1234

<file raw data>
```

or

```html
<form action="/v1/upload" method="post" enctype="multipart/form-data">
    <input type="file" name="file">
    <input type="submit" value="Upload">
</form>
```

## pkg/upload/upload.go

```golang
package upload

import (
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
)

type uploadFile struct {
    fileName     string
    file         io.ReadSeekCloser
    size         int64
    deletingPath string // the path of the file to be deleted
}

var _ io.ReadSeekCloser = (*uploadFile)(nil)

func (f *uploadFile) Read(p []byte) (n int, err error) {
    if f.file == nil {
        return 0, errors.New("Can't read from a closed file")
    }
    return f.file.Read(p)
}

func (f *uploadFile) Seek(offset int64, whence int) (int64, error) {
    if f.file == nil {
        return 0, errors.New("Can't seek in a closed file")
    }
    return f.Seek(offset, whence)
}

func (f *uploadFile) Close() error {
    if f.file != nil {
        err := f.file.Close()
        if f.deletingPath != "" { // delete the temporary file
            _ = os.Remove(f.deletingPath)
        }
        return err
    }

    f.file = nil
    return nil
}

func (f *uploadFile) Name() string {
    return f.fileName
}

func (f *uploadFile) Size() int64 {
    return f.size
}

// GetFileFromRequest get the file from the request
// 1. read from the Body if the request is a "multipart/form-data"
// 2. read body as file, read file name from "?{fieldName}=fileName"
func GetFileFromRequest(r *http.Request, fieldName string, maxSize int64) (*uploadFile, error) {
    if r.Body == nil || r.ContentLength == 0 {
        return nil, errors.New("上传文件不能为空")
    }

    defer r.Body.Close()

    // Multipart form
    if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
        // ParseMultipartForm parses a request body as multipart/form-data
        if err := r.ParseMultipartForm(maxSize); err != nil {
            return nil, errors.New("Please upload a valid file with the 'multipart/form-data' request")
        }

        file, handler, err := r.FormFile(fieldName)

        if err != nil {
            return nil, err
        }

        return &uploadFile{
            fileName: handler.Filename,
            file:     file,
            size:     handler.Size,
        }, nil
    }

    if r.ContentLength > maxSize {
        return nil, fmt.Errorf("the file size exceeds the maximum size: %.2fMB", float64(maxSize)/1024./1024.)
    }

    // Body as file
    fileName := r.URL.Query().Get(fieldName)
    file, err := os.CreateTemp(os.TempDir(), "kratos-upload-*")
    if err != nil {
        return nil, err
    }
    contentLength := r.ContentLength
    if contentLength, err = io.Copy(file, r.Body); err != nil {
        return nil, err
    }

    // Seek to the beginning of the file
    _, _ = file.Seek(0, 0)

    return &uploadFile{
        fileName:     fileName,
        file:         file,
        size:         contentLength,
        deletingPath: file.Name(),
    }, nil
}

```