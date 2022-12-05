// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"sync"
)

type FileStorage interface {
	Save(fileHeader *multipart.FileHeader) (string, error)
}

var _ FileStorage = &LocalFileStorage{}

type LocalFileStorage struct {
	MediaRoot string
	lock      sync.RWMutex
}

// SetMediaRouter Set file base save path for LocalFileStorage
func (l *LocalFileStorage) SetMediaRouter(mediaRoot string) {
	l.MediaRoot = mediaRoot
}

// Save implement FileStorage
// Uploads the form file to local
func (l *LocalFileStorage) Save(fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := l.getAlternativeName(fileHeader.Filename)
	if err != nil {
		return "", err
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, src)
	if err != nil {
		return "", err
	}
	return dst, nil
}

// Return an alternative filename, by adding an underscore and a random 7
// character alphanumeric string (before the file extension, if one exists
// to the filename
func (l *LocalFileStorage) getAlternativeName(filename string) (string, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	for {
		dst := path.Join(l.MediaRoot, filename)
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
