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

package grouper

import (
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/log"
	"github.com/GoogleCloudPlatform/khi/pkg/testutil/testlog"

	_ "github.com/GoogleCloudPlatform/khi/internal/testflags"
)

func TestAllDependentLogGrouper(t *testing.T) {
	tests := []struct {
		name     string
		logs     []*log.Log
		wantKeys map[string]struct{}
	}{
		{
			name:     "empty logs",
			logs:     []*log.Log{},
			wantKeys: map[string]struct{}{},
		},
		{
			name: "simple case",
			logs: []*log.Log{
				testlog.NewEmptyLogWithID("id1"),
				testlog.NewEmptyLogWithID("id2"),
				testlog.NewEmptyLogWithID("id3"),
			},
			wantKeys: map[string]struct{}{
				"": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := AllDependentLogGrouper
			got := g.Group(tt.logs)
			if len(got) != len(tt.wantKeys) {
				t.Errorf("Key length mismatch")
			}
			for wantKey := range tt.wantKeys {
				_, found := got[wantKey]
				if !found {
					t.Errorf("key %s was not found in the result", wantKey)
				}
			}
		})
	}
}

func TestAllIndependentLogGrouper(t *testing.T) {
	tests := []struct {
		name     string
		logs     []*log.Log
		wantKeys map[string]struct{}
	}{
		{
			name:     "empty logs",
			logs:     []*log.Log{},
			wantKeys: map[string]struct{}{},
		},
		{
			name: "simple case",
			logs: []*log.Log{
				testlog.NewEmptyLogWithID("id1"),
				testlog.NewEmptyLogWithID("id2"),
				testlog.NewEmptyLogWithID("id3"),
			},
			wantKeys: map[string]struct{}{
				"id1": {},
				"id2": {},
				"id3": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := AllIndependentLogGrouper
			got := g.Group(tt.logs)
			if len(got) != len(tt.wantKeys) {
				t.Errorf("Key length mismatch")
			}
			for wantKey := range tt.wantKeys {
				_, found := got[wantKey]
				if !found {
					t.Errorf("key %s was not found in the result", wantKey)
				}
			}
		})
	}
}

func TestSingleStringFieldKeyLogGrouper(t *testing.T) {
	tests := []struct {
		name     string
		logs     []*log.Log
		wantKeys map[string]struct{}
	}{
		{
			name:     "empty logs",
			logs:     []*log.Log{},
			wantKeys: map[string]struct{}{},
		},
		{
			name: "multiple logs",
			logs: []*log.Log{
				testlog.MustLogFromYAML("textPayload: log message 1\nkey: groupA"),
				testlog.MustLogFromYAML("textPayload: log message 2\nkey: groupB"),
				testlog.MustLogFromYAML("textPayload: log message 3\nkey: groupA"),
			},
			wantKeys: map[string]struct{}{
				"groupA": {},
				"groupB": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewSingleStringFieldKeyLogGrouper("key")
			got := g.Group(tt.logs)
			if len(got) != len(tt.wantKeys) {
				t.Errorf("Key length mismatch")
			}
			for wantKey := range tt.wantKeys {
				_, found := got[wantKey]
				if !found {
					t.Errorf("key %s was not found in the result", wantKey)
				}
			}
		})
	}
}
