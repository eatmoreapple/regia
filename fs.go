// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const randomStringChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var maxRandomStringCharsLength = len(randomStringChars)

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

	// check path if exists
	if len(l.mediaRoot) > 0 {
		_, err = os.Stat(l.mediaRoot)
		if err != nil {
			if os.IsNotExist(err) {
				if err = os.Mkdir(l.mediaRoot, 0666); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	dst, err := l.getAlternativeName(fileHeader.Filename)
	if err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// Return an alternative filename, by adding an underscore and a random 7
// character alphanumeric string (before the file extension, if one exists
// to the filename
func (l LocalFileStorage) getAlternativeName(filename string) (string, error) {
	for {
		dst := path.Join(l.mediaRoot, filename)
		exist, err := fileExists(dst)
		if err != nil {
			return "", err
		}
		if !exist {
			return dst, err
		}
		filename = getRandomString(7) + filename
	}
}

// check file if exists
func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Return a securely generated random string
func getRandomString(length int) string {
	if length > maxRandomStringCharsLength {
		length = maxRandomStringCharsLength
	}
	var builder strings.Builder
	rand.Seed(time.Now().Unix())
	for i := 0; i < length; i++ {
		index := rand.Intn(maxRandomStringCharsLength)
		builder.WriteString(string(randomStringChars[index]))
	}
	return builder.String()
}
