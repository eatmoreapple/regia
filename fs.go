// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

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
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}
	return http.DetectContentType(buffer), nil
}

type FileStorage interface {
	Save(fileHeader *multipart.FileHeader) error
}

var _ FileStorage = LocalFileStorage{}

type LocalFileStorage struct {
	mediaRoot string
}

// SetMediaRouter Set file base save path for LocalFileStorage
func (l *LocalFileStorage) SetMediaRouter(mediaRoot string) {
	l.mediaRoot = mediaRoot
}

// Save implement FileStorage
// Uploads the form file to local
func (l LocalFileStorage) Save(fileHeader *multipart.FileHeader) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst := path.Join(l.mediaRoot, fileHeader.Filename)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
