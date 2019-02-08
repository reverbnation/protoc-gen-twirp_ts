package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	_ "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/golang/protobuf/ptypes/struct"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/require"
)

func testdataFileDescriptors(name string) []*descriptor.FileDescriptorProto {
	data, err := ioutil.ReadFile("testdata/gen/" + name + ".pb")
	if err != nil {
		panic(err)
	}
	v := &descriptor.FileDescriptorSet{}
	if err := proto.Unmarshal(data, v); err != nil {
		panic(err)
	}
	return v.File
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name        string
		descriptors string
		want        *plugin.CodeGeneratorResponse
		wantErr     bool
	}{
		{
			name:    "example",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &plugin.CodeGeneratorRequest{
				ProtoFile: testdataFileDescriptors(tt.name),
			}
			got, err := generate(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			saveOutput(t, got)

			for _, f := range got.File {
				want, err := ioutil.ReadFile(testOutputFilename(f.GetName()))
				require.NoError(t, err)

				if string(want) != f.GetContent() {
					diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
						A:        difflib.SplitLines(string(want)),
						B:        difflib.SplitLines(f.GetContent()),
						FromFile: "Want",
						FromDate: "",
						ToFile:   "Got",
						ToDate:   "",
						Context:  1,
					})
					require.NoError(t, err)
					t.Errorf("\n%s", diff)
				}
			}
		})
	}
}

func saveOutput(t *testing.T, res *plugin.CodeGeneratorResponse) {
	require.NoError(t, os.RemoveAll("testdata/output/"))
	for _, f := range res.File {
		outfile := testOutputFilename(f.GetName())
		require.NoError(t, os.MkdirAll(filepath.Dir(outfile), 0755))
		require.NoError(t, ioutil.WriteFile(outfile, []byte(f.GetContent()), 0644))
	}
}

func testOutputFilename(filename string) string {
	return fmt.Sprintf("testdata/output/%s", filename)
}

// Regenerate the input descriptor:
//go:generate protoc -I testdata/protos --include_imports --descriptor_set_out=testdata/gen/example.pb testdata/protos/example.proto
