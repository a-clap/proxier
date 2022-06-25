package config

import (
	"testing"
)

func TestNew(t *testing.T) {

	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
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
				},
				"files":{
				}
			}`)},
			wantErr: false,
		},
		{
			name:    "Not a json",
			args:    args{[]byte(`I am not json file`)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestConfig_Get(t *testing.T) {
	type args struct {
		buf  []byte
		keys []string
	}
	type wants struct {
		err   []bool
		value []string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "Just settings",
			args: args{
				buf: []byte(`
				{
				"settings": {
					"user": "adam",
					"password": "pwd",
					"http_proxy_server": "192.168.0.1",
					"port": "123"
					},
					"files":{
					}
				}`),
				keys: []string{"user", "password", "http_proxy_server", "port", "ports"},
			},
			wants: wants{
				err:   []bool{false, false, false, false, true},
				value: []string{"adam", "pwd", "192.168.0.1", "123", ""},
			},
		},
		{
			name: "Settings with single variable",
			args: args{
				buf: []byte(`
				{
				"settings": {
					"user": "adam",
					"password": "pwd",
					"http_proxy_server": "192.168.0.1",
					"port": "123",
					"http_proxy": "${user}"
					},
				"files":{
					}
				}`),
				keys: []string{"http_proxy"},
			},
			wants: wants{
				err:   []bool{false},
				value: []string{"adam"},
			},
		},
		{
			name: "Settings with variable - 2 same keys",
			args: args{
				buf: []byte(`
				{
				"settings": {
					"user": "adam",
					"password": "pwd",
					"http_proxy_server": "192.168.0.1",
					"port": "123",
					"http_proxy": "${user}123${user}"
					},
				"files": {
					}
				}`),
				keys: []string{"http_proxy"},
			},
			wants: wants{
				err:   []bool{false},
				value: []string{"adam123adam"},
			},
		},
		{
			name: "Settings with variables - 2 different keys",
			args: args{
				buf: []byte(`
				{	
				"settings": {
					"user": "adam",
					"password": "pwd",
					"http_proxy_server": "192.168.0.1",
					"port": "123",
					"http_proxy": "${user}:${password}"
					},
				"files":{
					}
				}`),
				keys: []string{"http_proxy"},
			},
			wants: wants{
				err:   []bool{false},
				value: []string{"adam:pwd"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.args.buf)
			if err != nil {
				t.Errorf("Wrong buf = %v, err = %v", string(tt.args.buf), err)
				return
			}
			for i := 0; i < len(tt.args.keys); i++ {
				got, err := c.Get(tt.args.keys[i])
				if (err != nil) != tt.wants.err[i] {
					t.Errorf("Get() error = %v, err %v", err, tt.wants.err[i])
					return
				}
				if got != tt.wants.value[i] {
					t.Errorf("Get() got = %v, want %v", got, tt.wants.value[i])
				}
			}
		})
	}
}
