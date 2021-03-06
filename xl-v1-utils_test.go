/*
 * Minio Cloud Storage, (C) 2015, 2016 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import "testing"

// Test for reduceErrs.
func TestReduceErrs(t *testing.T) {
	// List all of all test cases to validate various cases of reduce errors.
	testCases := []struct {
		errs        []error
		ignoredErrs []error
		err         error
	}{
		// Validate if have reduced properly.
		{[]error{
			errDiskNotFound,
			errDiskNotFound,
			errDiskFull,
		}, []error{}, errDiskNotFound},
		// Validate if have no consensus.
		{[]error{
			errDiskFull,
			errDiskNotFound,
			nil, nil,
		}, []error{}, nil},
		// Validate if have consensus and errors ignored.
		{[]error{
			errVolumeNotFound,
			errVolumeNotFound,
			errVolumeNotFound,
			errDiskNotFound,
			errDiskNotFound,
		}, []error{errDiskNotFound}, errVolumeNotFound},
		{[]error{}, []error{}, nil},
	}
	// Validates list of all the testcases for returning valid errors.
	for i, testCase := range testCases {
		gotErr := reduceErrs(testCase.errs, testCase.ignoredErrs)
		if testCase.err != gotErr {
			t.Errorf("Test %d : expected %s, got %s", i+1, testCase.err, gotErr)
		}
	}
}

// Test for unionChecksums
func TestUnionChecksumInfos(t *testing.T) {
	cur := []checkSumInfo{
		{"part.1", "dummy", "cur-hash.1"},
		{"part.2", "dummy", "cur-hash.2"},
		{"part.3", "dummy", "cur-hash.3"},
		{"part.4", "dummy", "cur-hash.4"},
		{"part.5", "dummy", "cur-hash.5"},
	}
	updated := []checkSumInfo{
		{"part.1", "dummy", "updated-hash.1"},
		{"part.2", "dummy", "updated-hash.2"},
		{"part.3", "dummy", "updated-hash.3"},
	}
	curPartcksum := cur[0] // part.1 is the current part being written

	// Verify that hash of current part being written must be from cur []checkSumInfo
	finalChecksums := unionChecksumInfos(cur, updated, curPartcksum.Name)
	for _, cksum := range finalChecksums {
		if cksum.Name == curPartcksum.Name && cksum.Hash != curPartcksum.Hash {
			t.Errorf("expected Hash = %s but received Hash = %s\n", curPartcksum.Hash, cksum.Hash)
		}
	}

	// Verify that all part checksums are present in the union and nothing more.
	// Map to store all unique part names
	allPartNames := make(map[string]struct{})
	// Insert part names from cur and updated []checkSumInfo
	for _, cksum := range cur {
		allPartNames[cksum.Name] = struct{}{}
	}
	for _, cksum := range updated {
		allPartNames[cksum.Name] = struct{}{}
	}
	// All parts must have an entry in the []checkSumInfo returned from unionChecksums
	for _, finalcksum := range finalChecksums {
		if _, ok := allPartNames[finalcksum.Name]; !ok {
			t.Errorf("expected to find %s but not present in the union, where current part is %s\n",
				finalcksum.Name, curPartcksum.Name)
		}
	}
	if len(finalChecksums) != len(allPartNames) {
		t.Error("Union of Checksums doesn't have same number of elements as unique parts in total")
	}
}
