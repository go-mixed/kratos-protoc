package streaming

import "io"

type IProtobuf interface {
}

type Streaming struct {
	contentType string
	reader      io.ReadCloser
}

func NewStreaming(contentType string, reader io.ReadCloser) *Streaming {
	return &Streaming{
		contentType: contentType,
		reader:      reader,
	}
}

func (s *Streaming) StreamContentType() string {
	return s.contentType
}

func (s *Streaming) StreamReadCloser() io.ReadCloser {
	return s.reader
}
