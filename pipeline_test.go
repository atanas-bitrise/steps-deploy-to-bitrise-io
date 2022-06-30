package main

import (
	"reflect"
	"testing"
)

func Test_parsePipelineIntermediateFiles(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "simple test",
			s:       "$BITRISE_IPA_PATH:BITRISE_IPA_PATH",
			want:    map[string]string{"$BITRISE_IPA_PATH": "BITRISE_IPA_PATH"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePipelineIntermediateFiles(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePipelineIntermediateFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePipelineIntermediateFiles() got = %v, want %v", got, tt.want)
			}
		})
	}
}
