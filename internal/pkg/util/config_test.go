// Package util contains common functionality for "utilities" required by the
// bot and feed poller.
//
// This file is for processing the application's configuration from a YAML file.
package util

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func Test_loadConfigData(t *testing.T) {
	type args struct {
		params []string
	}

	tests := []struct {
		name       string
		args       args
		want       []byte
		wantErr    bool
		errMessage string
	}{
		{
			name:       "no filename provided and no config file exists",
			args:       args{},
			want:       nil,
			wantErr:    true,
			errMessage: "config.yaml doesn't exist in the global settings directory",
		},
		{
			name:       "filename provided and no config file exists",
			args:       args{params: []string{"foo"}},
			want:       nil,
			wantErr:    true,
			errMessage: "open foo: no such file or directory",
		},
		{
			name: "filename provided and config file exists",
			args: args{params: []string{"../../../test/config.yaml"}},
			want: func() []byte {
				data, _ := ioutil.ReadFile("../../../test/config.yaml")

				return data
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadConfigData(tt.args.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfigData() error = %v, wantErr %v (message: %s)", err, tt.wantErr, err.Error())
				return
			}

			if err != nil && err.Error() != tt.errMessage {
				t.Errorf("loadConfigData() message = %s, errMessage = %s", err.Error(), tt.errMessage)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfigData() = %v, want %v", got, tt.want)
			}
		})
	}
}
