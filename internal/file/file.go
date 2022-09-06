package file

import (
	"errors"
	"fmt"
	"github.com/a-clap/logger"
	"io"
	"os"
)

var Logger logger.Logger = logger.NewNop()

type File interface {
	io.Reader
	io.WriterAt
	io.Closer
	Name() string
}

type FS interface {
	Stat(name string) (os.FileInfo, error)
	Create(name string) (File, error)
	Open(name string) (File, error)
}

type Handler struct {
	fs FS
}

var (
	ErrIsDirectory = errors.New("is directory")
)

func New(fs FS) (Handler, error) {
	return Handler{fs}, nil
}

func (h *Handler) Backup(src, dst string) error {
	srcInfo, err := h.fs.Stat(src)
	if err != nil {
		return err
	}
	if srcInfo.IsDir() {
		return ErrIsDirectory
	}

	_, err = h.fs.Stat(dst)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	backupFile, err := h.fs.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating file %v", err)
	}

	defer func(file File) {
		err := file.Close()
		if err != nil {
			Logger.Errorf("failed on close file %s", file.Name())
		}
	}(backupFile)

	srcFile, err := h.fs.Open(src)
	if err != nil {
		return fmt.Errorf("error opening srcFile %s %v", src, err)
	}
	defer func(file File) {
		err := file.Close()
		if err != nil {
			Logger.Errorf("failed on close file %s", file.Name())
		}
	}(srcFile)

	srcBuf := make([]byte, srcInfo.Size())
	b, err := srcFile.Read(srcBuf)
	if err != nil {
		return fmt.Errorf("error reading file %s %v", src, err)
	}

	var pos int64 = 0
	for pos < int64(b) {
		written, err := backupFile.WriteAt(srcBuf, pos)
		if err != nil {
			return fmt.Errorf("error writing file %s, at pos %d", backupFile.Name(), pos)
		}
		pos += int64(written)
	}

	return nil
}
