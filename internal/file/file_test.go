package file_test

import (
	"errors"
	"github.com/a-clap/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"io/fs"
	"io/ioutil"
	"os"
	"proxier/internal/file"
	"syscall"
	"testing"
)

type mockFs struct {
	fs afero.Fs
}

func (m *mockFs) Stat(name string) (os.FileInfo, error) {
	return m.fs.Stat(name)
}

func (m *mockFs) Open(name string) (file.File, error) {
	return m.fs.Open(name)
}

func (m *mockFs) Create(name string) (file.File, error) {
	return m.fs.Create(name)
}

var _ file.FS = &mockFs{afero.NewOsFs()}

func TestHandler_Backup(t *testing.T) {
	logger.Init(logger.NewDefaultZap(zapcore.DebugLevel))

	type fields struct {
		fs mockFs
	}
	type dir struct {
		name string
		perm os.FileMode
	}
	type tmpfile struct {
		name string
		data []byte
	}
	type prepare struct {
		dir  dir
		file tmpfile
	}
	type args struct {
		src, dest string
	}
	tests := []struct {
		name        string
		fields      fields
		prepare     []prepare
		args        args
		wantErr     bool
		wantErrType error
		wantData    []byte
	}{
		{
			name:        "Source file doesn't exists",
			fields:      fields{mockFs{afero.NewOsFs()}},
			prepare:     nil,
			args:        args{"not_exist.txt", "backup.txt"},
			wantErr:     true,
			wantErrType: &fs.PathError{Op: "stat", Path: "not_exist.txt", Err: syscall.ENOENT},
		},
		{
			name:   "Source file is dir",
			fields: fields{mockFs{afero.NewOsFs()}},
			prepare: []prepare{
				{
					dir: dir{
						name: "dir.txt",
						perm: 0777,
					},
					file: tmpfile{
						name: "",
						data: nil,
					},
				},
			},
			args:        args{"dir.txt", "backup.txt"},
			wantErr:     true,
			wantErrType: file.ErrIsDirectory,
		},
		{
			name:   "No permission to src file",
			fields: fields{mockFs{afero.NewOsFs()}},
			prepare: []prepare{
				{
					dir: dir{
						name: "no_permission",
						perm: 0000,
					},
					file: tmpfile{
						name: "no_permission/src.txt",
						data: nil,
					},
				},
			},
			args:        args{"no_permission/src.txt", "backup.txt"},
			wantErr:     true,
			wantErrType: &fs.PathError{Op: "stat", Path: "no_permission/src.txt", Err: syscall.EACCES},
		},
		{
			name:   "No permission to dst file",
			fields: fields{mockFs{afero.NewOsFs()}},
			prepare: []prepare{
				{
					dir: dir{
						name: "permission",
						perm: 0777,
					},
					file: tmpfile{
						name: "permission/src.txt",
						data: nil,
					},
				},
				{
					dir: dir{
						name: "no_permission",
						perm: 0000,
					},
					file: tmpfile{
						name: "no_permission/backup.txt",
						data: nil,
					},
				},
			},
			args:        args{"permission/src.txt", "no_permission/backup.txt"},
			wantErr:     true,
			wantErrType: &fs.PathError{Op: "stat", Path: "no_permission/backup.txt", Err: syscall.EACCES},
		},
		{
			name:   "Successful copy",
			fields: fields{mockFs{afero.NewOsFs()}},
			prepare: []prepare{
				{
					dir: dir{
						name: "permission",
						perm: 0777,
					},
					file: tmpfile{
						name: "permission/src.txt",
						data: []byte("some useful data"),
					},
				},
				{
					dir: dir{
						name: "",
						perm: 0000,
					},
					file: tmpfile{
						name: "permission/backup.txt",
						data: nil,
					},
				},
			},
			args:     args{"permission/src.txt", "permission/backup.txt"},
			wantErr:  false,
			wantData: []byte("some useful data"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := file.New(&tt.fields.fs)
			require.Nil(t, err)

			for _, elem := range tt.prepare {
				if len(elem.dir.name) > 0 {
					if err := tt.fields.fs.fs.Mkdir(elem.dir.name, elem.dir.perm); err != nil {
						panic(err)
					}
				}

				if len(elem.file.name) > 0 {
					f, err := tt.fields.fs.fs.Create(elem.file.name)
					if err != nil {
						if !errors.Is(err, os.ErrPermission) {
							panic(err)
						}
					}

					if elem.file.data != nil {
						n, err := f.Write(elem.file.data)
						if n != len(elem.file.data) || err != nil {
							panic(err)
						}
						_ = f.Close()
					}
				}
			}

			err = h.Backup(tt.args.src, tt.args.dest)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErrType, err)
			} else {
				assert.Nil(t, err)
				f, err := tt.fields.fs.Open(tt.args.dest)
				data, err := ioutil.ReadAll(f)
				assert.Nil(t, err)
				assert.Equal(t, tt.wantData, data)
			}

			// Cleanup files
			for _, elem := range tt.prepare {
				if len(elem.file.name) > 0 {
					_ = tt.fields.fs.fs.Remove(elem.file.name)
				}
				if len(elem.dir.name) > 0 {
					_ = tt.fields.fs.fs.RemoveAll(elem.dir.name)
				}
			}

		})
	}
}
