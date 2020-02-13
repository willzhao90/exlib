package blob

import (
	"io"
)

type BlobStore interface {
	PullObjectURL(string, string) (string, error)
	PutObjectURL(string, string) (string, error)
	PullObject(string, string, io.WriterAt) error
	PutObject(string, string, io.Reader) error
}
