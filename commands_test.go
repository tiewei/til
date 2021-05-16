/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import "testing"

func TestSplitArgs(t *testing.T) {
	testCases := map[string]struct {
		n   int
		in  []string
		out [2][]string
	}{
		"n is less than number of arguments": {
			n:   2,
			in:  []string{"1", "2", "3", "4"},
			out: [2][]string{{"1", "2"}, {"3", "4"}},
		},
		"n is greater than number of arguments": {
			n:   4,
			in:  []string{"1", "2", "3"},
			out: [2][]string{{"1", "2", "3"}, nil},
		},
		"n is zero": {
			n:   0,
			in:  []string{"1", "2"},
			out: [2][]string{nil, {"1", "2"}},
		},
		"no argument": {
			n:   1,
			in:  nil,
			out: [2][]string{nil, nil},
		},
		"arguments start with a flag": {
			n:   2,
			in:  []string{"-1", "2", "3", "4"},
			out: [2][]string{nil, {"-1", "2", "3", "4"}},
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			pos, flags := splitArgs(tc.n, tc.in)

			expectPos, expectFlags := tc.out[0], tc.out[1]

			if !equalSlices(expectPos, pos) {
				t.Errorf("Expected positional to equal %q, got %q", expectPos, pos)
			}
			if !equalSlices(expectFlags, flags) {
				t.Errorf("Expected flags to equal %q, got %q", expectFlags, flags)
			}
		})
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// by using a range, we implicitly consider [] and nil as equal
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
