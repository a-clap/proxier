package config_test

import (
	"github.com/stretchr/testify/require"
	"proxier/internal/config"
	"reflect"
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
				"files":[]
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
			_, err := config.New(tt.args.buf)
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
					"files":[]
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
				"files":[]
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
				"files": []
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
				"files":[]
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
			c, err := config.New(tt.args.buf)
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

func TestConfig_GetFiles(t *testing.T) {
	type args struct {
		buf []byte
	}
	type wants struct {
		files []config.File
	}
	tests := []struct {
		name string
		args args
		want wants
	}{
		{
			name: "No file",
			args: args{
				buf: []byte(`{
				"settings": {
					"user": "adam"
					},
				"files":[]
				}`),
			},
			want: wants{files: nil},
		},
		{
			name: "Single file without variables",
			args: args{
				buf: []byte(`{
				"settings": {
					"user": "adam"
					},
				"files": [
				{
					"name": "first_file",
					"append": ["single_line"],
					"remove": ["single_line"]
				}
				]}`),
			},
			want: wants{files: []config.File{
				{
					Name:   "first_file",
					Append: []string{"single_line"},
					Remove: []string{"single_line"},
				},
			}},
		},
		{
			name: "Single file with variables",
			args: args{
				buf: []byte(`{
				"settings": {
					"user": "first",
					"password": "second",
					"proxy": "third",	
					"random_key": "random_value"
					},
				"files": [
				{
					"name": "first_file",
					"append": ["${user}", "${password}"],
					"remove": ["${proxy}", "${random_key}"]
				}
				]}`),
			},
			want: wants{files: []config.File{
				{
					Name:   "first_file",
					Append: []string{"first", "second"},
					Remove: []string{"third", "random_value"},
				},
			}},
		},
		{
			name: "Multiple files with variables",
			args: args{
				buf: []byte(`{
				"settings": {
					"user": "first",
					"password": "second",
					"proxy": "third",	
					"random_key": "random_value"
					},
				"files": [
				{
					"name": "first_file",
					"append": ["${user}", "${password}"],
					"remove": ["${proxy}", "${random_key}"]
				},
				{
					"name": "second_file",
					"append": ["value", "2value"],
					"remove": ["${proxy}", "${user}:${password}:${proxy}"]
				}
				]}`),
			},
			want: wants{files: []config.File{
				{
					Name:   "first_file",
					Append: []string{"first", "second"},
					Remove: []string{"third", "random_value"},
				},
				{
					Name:   "second_file",
					Append: []string{"value", "2value"},
					Remove: []string{"third", "first:second:third"},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := config.New(tt.args.buf)
			if err != nil {
				t.Errorf("Wrong buf = %v, err = %v", string(tt.args.buf), err)
				return
			}
			if got := c.GetFiles(); !reflect.DeepEqual(got, tt.want.files) {
				t.Errorf("GetFiles() = %v, want %v", got, tt.want.files)
			}
		})
	}
}

func TestTemplate(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "Basic config",
			want: []byte(`
			{
				"settings":{
					"user": "user",
					"password": "password",
					"server": "192.168.0.100",
					"port": "80",
					"http_proxy": "\"http://${user}:${password}@${server}:${port}\"",
					"https_proxy": "\"https://${user}:${password}@${server}:${port}\""
				},
				"files":[
				{
					"name": "/etc/environment",
					"append": [
						"HTTP_PROXY=${http_proxy}"	
					],
					"remove": [
						"HTTP_PROXY"
					]
				},
				{	
					"name": "/etc/apt/apt.conf.d/proxy.conf",
					"append": [
						"Acquire::http::proxy ${http_proxy}",
						"Acquire::https::proxy ${http_proxy}"
					],
					"remove": [
						"Acquire"
					]
				}
				]
			}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.Template()
			require.JSONEq(t, string(got), string(tt.want))
		})
	}
}
