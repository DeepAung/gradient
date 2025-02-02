package storer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

var (
	ErrUploadExistingDest      = errors.New("upload into existing destination")
	ErrDeleteNotExistingDest   = errors.New("delete at not existing destination")
	ErrDownloadNotExistingDest = errors.New("download from not existing destination")
)

type gcpStorer struct {
	gcpBucketName string
}

func NewGcpStorer(gcpBucketName string) Storer {
	return &gcpStorer{
		gcpBucketName: gcpBucketName,
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

	o := client.Bucket(s.gcpBucketName).Object(dest)

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

	o := client.Bucket(s.gcpBucketName).Object(dest)

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

func (s *gcpStorer) DeleteFolder(dir string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	it := client.Bucket(s.gcpBucketName).Objects(ctx, &storage.Query{
		Prefix: dir,
	})

	maxGoroutines := 10
	sem := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	for {
		attr, err := it.Next()
		if err == iterator.Done {
			break
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(name string) {
			s.Delete(name)
			wg.Done()
			<-sem
		}(attr.Name)
	}
	wg.Wait()

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

	rc, err := client.Bucket(s.gcpBucketName).Object(dest).NewReader(ctx)
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

func (s *gcpStorer) Download(remoteDest string, localDest string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	f, err := os.Create(localDest)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	rc, err := client.Bucket(s.gcpBucketName).Object(remoteDest).NewReader(ctx)
	if err != nil {
		if err.Error() == storage.ErrObjectNotExist.Error() {
			return ErrDownloadNotExistingDest
		}
		return fmt.Errorf("Object(%q).NewReader: %w", remoteDest, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}

	fmt.Printf("Blob %v downloaded.\n", remoteDest)

	return nil
}

// TODO: check error handling
func (s *gcpStorer) DownloadFolder(remoteDir string, localDir string) (int, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	it := client.Bucket(s.gcpBucketName).Objects(ctx, &storage.Query{
		Prefix: remoteDir,
	})

	maxGoroutines := 10
	sem := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup
	_, cancel2 := context.WithCancelCause(context.Background())

	counter := 0
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:

		}
		attr, err := it.Next()
		if err == iterator.Done {
			break
		}
		counter++

		wg.Add(1)
		sem <- struct{}{}
		go func(name string) {
			err = s.Download(remoteDir+"/"+name, localDir+"/"+name)
			if err != nil {
				cancel2(err)
			}
			wg.Done()
			<-sem
		}(attr.Name)
	}
	wg.Wait()

	return counter, nil
}
