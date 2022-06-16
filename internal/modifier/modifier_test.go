package modifier_test

import (
	"proxier/internal/modifier"
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
			m := modifier.New(tt.args.buf)
			if got := m.Get(); !reflect.DeepEqual(got, tt.args.buf) {
				t.Errorf("New() = %v, want %v", got, tt.args.buf)
			}
		})
	}
}
