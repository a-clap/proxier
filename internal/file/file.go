package file

import (
	"errors"
	"fmt"
	"os"
)

type FS interface {
	Stat(name string) (os.FileInfo, error)
	Create(name string) (*os.File, error)
	Open(name string) (*os.File, error)
}

type Handler struct {
	fs FS
}

func New(fs FS) (Handler, error) {
	return Handler{fs}, nil
}

func (h *Handler) Backup(src, dst string) error {
	srcInfo, err := h.fs.Stat(src)
	if err != nil {
		return err
	}
	if srcInfo.IsDir() {
		return fmt.Errorf(

			"%s is directory", srcInfo.Name())
	}

	_, err = h.fs.Stat(dst)
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	backupFile, err := h.fs.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating file %v", err)
	}

	defer func(backupFile *os.File) {
		err := backupFile.Close()
		if err != nil {

		}
	}(backupFile)

	srcFile, err := h.fs.Open(src)
	if err != nil {
		return fmt.Errorf("error opening srcFile %s %v", src, err)
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {

		}
	}(srcFile)

	var srcBuf []byte
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

func (h *Handler) Test(s string) {
	if _, err := h.fs.Stat(s); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("file exists", s)
	}
}
