package regia

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

// GetFileContentType Get File contentType
// os.File impl multipart.File
func GetFileContentType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", err
	}
	file.Seek(0, 0)
	return http.DetectContentType(buffer), nil
}

type FileStorage interface {
	Save(fileHeader *multipart.FileHeader) error
}

type LocalFileStorage struct {
	mediaRouter string
}

func (l *LocalFileStorage) SetMediaRouter(mediaRouter string) {
	l.mediaRouter = mediaRouter
}

func (l LocalFileStorage) Save(fileHeader *multipart.FileHeader) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst := path.Join(l.mediaRouter, fileHeader.Filename)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
