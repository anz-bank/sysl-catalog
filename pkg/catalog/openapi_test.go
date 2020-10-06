package catalog

import (
	"testing"
)

func TestIsOpenAPIFile(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{"Handles empty string", "", false},
		{"Handles .sysl", "github.com/myorg/test/spec.sysl", false},
		{"Handles .oas.yaml", "github.com/myorg/test/spec.oas.yaml", true},
		{"Handles .yaml", "github.com/myorg/test/spec.yaml", true},
		{"Handles .yml", "github.com/myorg/test/spec.yml", true},
		{"Handles .json", "github.com/myorg/test/spec.json", true},
		{"Handles @develop", "github.com/myorg/test/spec.json@develop", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsOpenAPIFile(tt.arg)
			if got != tt.want {
				t.Errorf("IsOpenAPIFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildSpecURL(t *testing.T) {
	type args struct {
		filePath string
		version  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Simple", args{filePath: "./pkg/catalog/test/simple.yaml"}, "/pkg/catalog/test/simple.yaml", false},
		{"NoDot", args{filePath: "/pkg/catalog/test/simple.yaml"}, "/pkg/catalog/test/simple.yaml", false},
		{"AppendForwardSlash", args{filePath: "pkg/catalog/test/simple.yaml"}, "/pkg/catalog/test/simple.yaml", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildSpecURL(tt.args.filePath, tt.args.version)
			if got != tt.want {
				t.Errorf("BuildSpecURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
