package modifier_test

import (
	"proxier/internal/modifier"
	"proxier/pkg/logger"
	"reflect"
	"testing"
)

func TestNew_Get(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Simple one line buffer",
			args: args{buf: []byte("1234567")},
		},
		{
			name: "Two line buffer",
			args: args{buf: []byte("123\n456")},
		},
		{
			name: "Two line buffer with newline ending",
			args: args{buf: []byte("123\n456\n")},
		},
		{
			name: "Multiple line endings in buf",
			args: args{buf: []byte("123\n\n\n456\n\n\n789\n")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := modifier.New(tt.args.buf, logger.Dummy{})
			if got := m.Get(); !reflect.DeepEqual(got, tt.args.buf) {
				t.Errorf("New() = %v, want %v", got, tt.args.buf)
			}
		})
	}
}

func TestModifier_RemoveLines(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		pattern string
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantLinesRemoved int
		wantErr          bool
		wantBuf          fields
	}{
		{
			name:             "Remove one line with new line",
			fields:           fields{buf: []byte("keep this line\nremove this line\n")},
			args:             args{pattern: `remove this line`},
			wantLinesRemoved: 1,
			wantErr:          false,
			wantBuf:          fields{buf: []byte("keep this line\n")},
		},
		{
			name:             "Remove two lines",
			fields:           fields{buf: []byte("keep this line\nremove this line\nremove this line\n")},
			args:             args{pattern: `remove this line`},
			wantLinesRemoved: 2,
			wantErr:          false,
			wantBuf:          fields{buf: []byte("keep this line\n")},
		},
		{
			name:             "Remove many lines",
			fields:           fields{buf: []byte("keep this line\nremove this line\nremove this line\nremove this line\nremove this line\nremove this line\nremove this line\nremove this line\n")},
			args:             args{pattern: `remove this line`},
			wantLinesRemoved: 7,
			wantErr:          false,
			wantBuf:          fields{buf: []byte("keep this line\n")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := modifier.New(tt.fields.buf, logger.Dummy{})
			gotLinesRemoved, err := m.RemoveLines(tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLinesRemoved != tt.wantLinesRemoved {
				t.Errorf("RemoveLines() gotLinesRemoved = %v, want %v", gotLinesRemoved, tt.wantLinesRemoved)
			}
			if got := m.Get(); !reflect.DeepEqual(got, tt.wantBuf.buf) {
				t.Errorf("RemoveLines() gotBuf = %s\n, want %s\n", string(got), string(tt.wantBuf.buf))
			}
		})
	}
}

func TestModifier_AppendLines(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		lines []string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantLinesAppended int
		wantBuf           fields
	}{
		{
			name:              "",
			fields:            fields{},
			args:              args{},
			wantLinesAppended: 0,
			wantBuf:           fields{buf: []byte("")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := modifier.New(tt.fields.buf, logger.Dummy{})
			if gotLinesAppended := m.AppendLines(tt.args.lines); gotLinesAppended != tt.wantLinesAppended {
				t.Errorf("AppendLines() = %v, want %v", gotLinesAppended, tt.wantLinesAppended)
			}
			if got := m.Get(); !reflect.DeepEqual(got, tt.wantBuf.buf) {
				t.Errorf("RemoveLines() gotBuf = %s\n, want %s\n", string(got), string(tt.wantBuf.buf))
			}
		})
	}
}
