package file_test

import (
	"github.com/spf13/afero"
	"log"
	"os"
	"proxier/internal/file"
	"reflect"
	"testing"
)

type aferoInterface struct {
	af afero.Fs
}

func (a aferoInterface) Stat(name string) (os.FileInfo, error) {
	return a.af.Stat(name)
}
func (a aferoInterface) Create(name string) (*os.File, error) {
	created, err := a.af.Create(name)
	return created.(*os.File), err
}
func (a aferoInterface) Open(name string) (*os.File, error) {
	opened, err := a.af.Open(name)
	return opened.(*os.File), err
}

func TestHandler_Backup(t *testing.T) {

	inf := aferoInterface{af: afero.NewOsFs()}

	// Create directory and files
	_ = inf.af.Mkdir("dir.txt", 0755)
	_ = inf.af.Mkdir("no_permission", 0000)
	_ = inf.af.Mkdir("permission", 0777)
	_, _ = inf.af.Create("no_permission/src.txt")

	// Delete previous files
	_ = inf.af.Remove("permission/src.txt")
	_ = inf.af.Remove("permission/backup.txt")

	// Create new files
	src, _ := inf.af.Create("permission/src.txt")
	data := []byte("very useful data\n")
	_, _ = src.Write(data)
	_ = src.Close()

	type fields struct {
		fs aferoInterface
	}
	type args struct {
		src, dest string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		wantData []byte
	}{
		{
			name:    "Source file doesn't exists",
			fields:  fields{inf},
			args:    args{"not_exist.txt", "backup.txt"},
			wantErr: true,
		},
		{
			name:    "Source file is dir",
			fields:  fields{inf},
			args:    args{"dir.txt", "backup.txt"},
			wantErr: true,
		},
		{
			name:    "No permission to src file",
			fields:  fields{inf},
			args:    args{"no_permission/src.txt", "backup.txt"},
			wantErr: true,
		},
		{
			name:    "No permission to dst file",
			fields:  fields{inf},
			args:    args{"permission/src.txt", "no_permission/backup.txt"},
			wantErr: true,
		},
		{
			name:     "Successful copy",
			fields:   fields{inf},
			args:     args{"permission/src.txt", "permission/backup.txt"},
			wantErr:  false,
			wantData: data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := file.New(&tt.fields.fs)
			if err != nil {
				log.Fatalln(err)
			}
			if err := h.Backup(tt.args.src, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("Backup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				backupData, err := afero.ReadFile(inf.af, tt.args.dest)
				if err != nil {
					t.Errorf("Not expected error %v", err)
				}
				if !reflect.DeepEqual(backupData, tt.wantData) {
					t.Errorf("Wrong backup %v, %v", backupData, tt.wantData)
				}
			}

		})
	}
}
