// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testlog

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	_ "github.com/GoogleCloudPlatform/khi/internal/testflags"
	"github.com/GoogleCloudPlatform/khi/pkg/common/structured"
)

func TestBaseYamlTestLogOpt(t *testing.T) {
	testCases := []struct {
		name        string
		inputYaml   string
		outputYaml  string
		expectError bool
	}{
		{
			name:      "basic valid yaml",
			inputYaml: `foo: bar`,
			outputYaml: `foo: bar
`,
			expectError: false,
		},
		{
			name:        "parses empty yaml as an empty map",
			inputYaml:   "",
			outputYaml:  "{}\n",
			expectError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tl := New(YAML(tc.inputYaml))
			reader, err := tl.BuildReader()
			if tc.expectError {
				if err == nil {
					t.Errorf("Expecting an error but no error returned.")
				}
			} else {
				yamlStr, err := reader.Serialize("", &structured.YAMLNodeSerializer{})
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if diff := cmp.Diff(string(yamlStr), tc.outputYaml); diff != "" {
					t.Errorf("Yaml string mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
