package storer

import (
	"io"
	"mime/multipart"
	"path"
)

const basePath = "https://storage.googleapis.com/gradient-bucket-dev"

type Storer interface {
	Upload(reader io.Reader, dest string, public bool) (FileRes, error)
	UploadMultipart(f *multipart.FileHeader, dest string, public bool) (FileRes, error)

	Delete(dest string) error
	DeleteFolder(dir string) error

	Download(remoteDest string, localDest string) error
	DownloadContent(dest string) (string, error)
	DownloadFolder(remoteDir string, localDir string) (count int, err error)
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
