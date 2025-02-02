package storer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"github.com/DeepAung/gradient/website-server/config"
)

var (
	ErrUploadExistingDest      = errors.New("upload into existing destination")
	ErrDeleteNotExistingDest   = errors.New("delete at not existing destination")
	ErrDownloadNotExistingDest = errors.New("download from not existing destination")
)

type gcpStorer struct {
	cfg *config.Config
}

func NewGcpStorer(cfg *config.Config) Storer {
	return &gcpStorer{
		cfg: cfg,
	}
}

func (s *gcpStorer) UploadMultipart(
	fileHeader *multipart.FileHeader,
	dest string,
	public bool,
) (FileRes, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return FileRes{}, err
	}
	defer file.Close()
	return s.Upload(file, dest, public)
}

func (s *gcpStorer) Upload(reader io.Reader, dest string, public bool) (FileRes, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return FileRes{}, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	o := client.Bucket(s.cfg.App.GcpBucketName).Object(dest)

	o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, reader); err != nil {
		return FileRes{}, fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		if err.Error() == `googleapi: Error 412: At least one of the pre-conditions you specified did not hold., conditionNotMet` {
			return FileRes{}, ErrUploadExistingDest
		}
		return FileRes{}, fmt.Errorf("Writer.Close: %w", err)
	}
	fmt.Printf("Blob %v uploaded.\n", dest)

	if public {
		err := o.ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
		if err != nil {
			return FileRes{}, fmt.Errorf("ACL.Set: %w", err)
		}
	}

	return NewFileResFromDest(dest), nil
}

func (s *gcpStorer) Delete(dest string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	o := client.Bucket(s.cfg.App.GcpBucketName).Object(dest)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		if err.Error() == storage.ErrObjectNotExist.Error() {
			return ErrDeleteNotExistingDest
		}
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", dest, err)
	}
	fmt.Printf("Blob %v deleted.\n", dest)
	return nil
}

func (s *gcpStorer) DownloadContent(dest string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	buf := bytes.NewBufferString("")

	rc, err := client.Bucket(s.cfg.App.GcpBucketName).Object(dest).NewReader(ctx)
	if err != nil {
		if err.Error() == storage.ErrObjectNotExist.Error() {
			return "", ErrDownloadNotExistingDest
		}
		return "", fmt.Errorf("Object(%q).NewReader: %w", dest, err)
	}
	defer rc.Close()

	if _, err := io.Copy(buf, rc); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}

	fmt.Printf("Blob %v downloaded.\n", dest)

	return buf.String(), nil
}
