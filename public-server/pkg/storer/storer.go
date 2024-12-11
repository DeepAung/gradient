package storer

import (
	"io"
	"path"
)

const basePath = "https://storage.googleapis.com/gradient-bucket-dev"

type Storer interface {
	Upload(reader io.Reader, dest string, public bool) (FileRes, error)
	Delete(dest string) error
}

// url      = https://storage.googleapis.com/gradient-bucket-dev/users/1/profile.jpg
// basePath = https://storage.googleapis.com/gradient-bucket-dev
// dest     = users/1/profile.jpg
// dir      = users/1
// filename = profile.jpg
type FileRes struct {
	Url      string
	BasePath string
	Dest     string
	Dir      string
	Filename string
}

func NewFileResFromDest(dest string) FileRes {
	filename := path.Base(dest)
	dir := path.Dir(dest)
	return FileRes{
		Url:      path.Join(basePath, dest),
		BasePath: basePath,
		Dest:     dest,
		Dir:      dir,
		Filename: filename,
	}
}
