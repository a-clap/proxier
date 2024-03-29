package proxier

import (
	"fmt"
	"github.com/a-clap/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"proxier/internal/config"
	"proxier/internal/file"
	"proxier/internal/modifier"
	"time"
)

var Logger logger.Logger = logger.NewNop()

type Proxier struct {
	cfg         *config.Config
	fileHandler file.Handler
}
type OsFs struct {
}

func (o *OsFs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
func (o *OsFs) Create(name string) (file.File, error) {
	return os.Create(name)
}
func (o *OsFs) Open(name string) (file.File, error) {
	return os.Open(name)
}

var _ file.FS = &OsFs{}

const CONFIG_FILE = "config.json"

func New() (*Proxier, error) {
	f, err := os.Open(CONFIG_FILE)
	if err != nil {
		return nil, fmt.Errorf("error opening config file %v = %v", CONFIG_FILE, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			Logger.Errorf("Error closing %s file = %v\n", CONFIG_FILE, err)
		}
	}(f)

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %v = %v", CONFIG_FILE, err)
	}

	cfg, err := config.New(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file %v = %v", CONFIG_FILE, err)
	}
	fileHandler, _ := file.New(&OsFs{})

	return &Proxier{cfg: cfg,
		fileHandler: fileHandler,
	}, nil
}

func (p *Proxier) Set(backup bool) error {
	if backup {
		if err := p.backup(); err != nil {
			return fmt.Errorf("error on backup %v", err)
		}
	}

	for _, currentFile := range p.cfg.GetFiles() {
		Logger.Infof("Appending lines to file \"%s\"", currentFile.Name)
		f, err := os.OpenFile(currentFile.Name, os.O_RDWR, os.ModeAppend)
		if err != nil {
			return fmt.Errorf("error opening file %v = %v", currentFile.Name, err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				Logger.Fatalf("Failed closing file %s, error = %v", currentFile.Name, err)
			}
		}(f)
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("error reading file %v = %v", currentFile.Name, err)
		}

		m := modifier.New(data)
		m.AppendLines(currentFile.Append)
		_, err = f.WriteAt(m.Get(), 0)
		if err != nil {
			return fmt.Errorf("failed writting file %v\n", err)
		}
	}

	return nil
}

func (p *Proxier) Unset(backup bool) error {
	if backup {
		if err := p.backup(); err != nil {
			return fmt.Errorf("error on backup %v", err)
		}
	}

	for _, currentFile := range p.cfg.GetFiles() {
		Logger.Infof("Removing lines from file \"%s\"", currentFile.Name)
		f, err := os.OpenFile(currentFile.Name, os.O_RDWR, os.ModeAppend)
		if err != nil {
			return fmt.Errorf("error opening file %v = %v", currentFile.Name, err)
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				Logger.Fatalf("Failed closing file %s, error = %v", currentFile.Name, err)
			}
		}(f)

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("error reading file %v = %v", currentFile.Name, err)
		}

		m := modifier.New(data)
		for _, line := range currentFile.Remove {
			_, err := m.RemoveLines(line)
			if err != nil {
				return err
			}
		}

		err = os.Truncate(currentFile.Name, 0)
		if err != nil {
			return err
		}

		_, err = f.WriteAt(m.Get(), 0)
		if err != nil {
			return fmt.Errorf("failed writting file %v\n", err)
		}
	}

	return nil
}

func (p *Proxier) backup() error {
	dstName := "_" + time.Now().Format("02_01_2006_15_04.0000")
	for _, backupFile := range p.cfg.GetFiles() {
		backupDir := filepath.Dir(backupFile.Name) + "/backup"
		err := os.Mkdir(backupDir, 0755)

		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create backup dir %v %v", backupDir, err)
		}

		fileName := filepath.Base(backupFile.Name) + dstName

		if err := p.fileHandler.Backup(backupFile.Name, backupDir+"/"+fileName); err != nil {
			Logger.Errorf("failed on doing backup of file %v = %v", backupFile.Name, err)
		}
	}
	return nil
}
