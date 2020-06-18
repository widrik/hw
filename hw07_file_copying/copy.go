package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	oldFile, err := os.Open(fromPath)

	if err != nil {
		return ErrUnsupportedFile
	}
	defer oldFile.Close()

	fileInfo, err := oldFile.Stat()

	if err != nil {
		return err
	}

	switch {
	case fileInfo.IsDir():
		return ErrUnsupportedFile
	case fileInfo.Size() < offset:
		return ErrOffsetExceedsFileSize
	default:
	}

	if _, err := oldFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	newFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	count := int(fileInfo.Size() - offset)
	bar := pb.StartNew(count)
	defer bar.Finish()

	barReader := bar.NewProxyReader(oldFile)
	_, err = io.CopyN(newFile, barReader, limit)
	if err != nil {
		return err
	}

	return nil
}
