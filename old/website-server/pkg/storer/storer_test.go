package storer

import (
	"testing"

	"github.com/DeepAung/gradient/website-server/pkg/asserts"
)

var defaultRes = FileRes{
	Url:      "https://storage.googleapis.com/gradient-bucket-dev/users/1/profile.jpg",
	BasePath: "https://storage.googleapis.com/gradient-bucket-dev",
	Dest:     "users/1/profile.jpg",
	Dir:      "users/1",
	Filename: "profile.jpg",
}

func TestNewFileResFromDest(t *testing.T) {
	res := NewFileResFromDest(defaultRes.Dest)
	asserts.Equal(t, "res", res, defaultRes)
}

func TestNewFileResFromUrl(t *testing.T) {
	res := NewFileResFromUrl(defaultRes.Url)
	asserts.Equal(t, "res", res, defaultRes)
}
