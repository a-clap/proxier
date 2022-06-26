package file_test

import (
	"github.com/spf13/afero"
	"log"
	"proxier/internal/file"
	"testing"
)

func TestHandler_Backup(t *testing.T) {
	testFs := afero.NewOsFs()

	type inf struct {
	}

	type fields struct {
		fs file.FS
	}
	type args struct {
		src, dest string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Do backup",
			fields:  fields{testFs},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := file.New(tt.fields.fs)
			if err != nil {
				log.Fatalln(err)
			}
			if err := h.Backup(tt.args.src, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("Backup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
