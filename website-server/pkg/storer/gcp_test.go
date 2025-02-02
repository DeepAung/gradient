package storer

import (
	"bytes"
	"testing"

	"github.com/DeepAung/gradient/website-server/config"
	"github.com/DeepAung/gradient/website-server/pkg/asserts"
)

var (
	cfg    *config.Config
	storer Storer

	privateDest  = "tests/private.txt"
	publicDest   = "tests/public.txt"
	notExistDest = "tests/not-exist.txt"
)

func init() {
	cfg = config.NewConfig("../../.env.dev")
	storer = NewGcpStorer(cfg)
}

func TestUpload(t *testing.T) {
	t.Run("private upload", func(t *testing.T) {
		reader := bytes.NewBufferString("hello world")
		res, err := storer.Upload(reader, privateDest, false)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "file result", res, NewFileResFromDest(privateDest))
	})

	t.Run("public upload", func(t *testing.T) {
		reader := bytes.NewBufferString("hello world")
		res, err := storer.Upload(reader, publicDest, true)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "file result", res, NewFileResFromDest(publicDest))
	})

	t.Run("upload into existing destination", func(t *testing.T) {
		reader := bytes.NewBuffer([]byte("hello world (updated)"))
		res, err := storer.Upload(reader, publicDest, false)
		asserts.EqualError(t, err, ErrUploadExistingDest)
		asserts.Equal(t, "file result", res, FileRes{})
	})
}

func TestDownloadContent(t *testing.T) {
	t.Run("public download", func(t *testing.T) {
		content, err := storer.DownloadContent(publicDest)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "content", content, "hello world")
	})

	t.Run("private download", func(t *testing.T) {
		content, err := storer.DownloadContent(privateDest)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "content", content, "hello world")
	})

	t.Run("destination not found", func(t *testing.T) {
		content, err := storer.DownloadContent(notExistDest)
		asserts.EqualError(t, err, ErrDownloadNotExistingDest)
		asserts.Equal(t, "content", content, "")
	})
}

func TestDelete(t *testing.T) {
	t.Run("destination not found", func(t *testing.T) {
		err := storer.Delete(notExistDest)
		asserts.EqualError(t, err, ErrDeleteNotExistingDest)
	})

	t.Run("delete private", func(t *testing.T) {
		err := storer.Delete(privateDest)
		asserts.EqualError(t, err, nil)
	})

	t.Run("delete public", func(t *testing.T) {
		err := storer.Delete(publicDest)
		asserts.EqualError(t, err, nil)
	})
}
