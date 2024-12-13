package storer

import (
	"io"
	"mime/multipart"
	"path"
)

const basePath = "https://storage.googleapis.com/gradient-bucket-dev"

type Storer interface {
	UploadMultipart(f *multipart.FileHeader, dest string, public bool) (FileRes, error)
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
		Url:      basePath + "/" + dest,
		BasePath: basePath,
		Dest:     dest,
		Dir:      dir,
		Filename: filename,
	}
}

func NewFileResFromUrl(url string) FileRes {
	// url = `{basePath}/{dest}`
	// exclude `{basePath}/`
	dest := url[len(basePath)+1:]
	return NewFileResFromDest(dest)
}
