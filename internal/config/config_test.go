package config

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {

	wantUser := "adam"
	wantPassword := "pwd"
	wantHttpProxy := "192.168.0.1"
	wantPort := "123"

	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "Correct config without variables",
			args: args{[]byte(`
			{
				"settings": {
					"user": "adam",
					"password": "pwd",
					"http_proxy_server": "192.168.0.1",
					"port": "123"
				}
			}`)},
			want: &Config{Settings: Settings{
				User:      &wantUser,
				Password:  &wantPassword,
				HttpProxy: &wantHttpProxy,
				Port:      &wantPort,
			}},
			wantErr: false,
		},
		{
			name: "Lack of user",
			args: args{[]byte(`
			{
				"settings": {
					"password": "pwd",
					"http_proxy_server": "192.168.0.1",
					"port": "123"
				}
			}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Lack of password",
			args: args{[]byte(`
			{
				"settings": {
					"user": "adam"
					"http_proxy_server": "192.168.0.1",
					"port": "123"
				}
			}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Lack of proxy server",
			args: args{[]byte(`
			{
				"settings": {
					"user": "adam",
					"password": "pwd",
					"port": "123"
				}
			}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Lack of port",
			args: args{[]byte(`
			{
				"settings": {
					"user": "adam",
					"password": "pwd",
					"http_proxy_server": "192.168.0.1"
				}
			}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Not a json",
			args:    args{[]byte(`I am not json file`)},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
