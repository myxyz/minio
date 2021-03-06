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

import (
	"io"
	"net/http"
)

// Validates location constraint in PutBucket request body.
// The location value in the request body should match the
// region configured at serverConfig, otherwise error is returned.
func isValidLocationConstraint(r *http.Request) (s3Error APIErrorCode) {
	serverRegion := serverConfig.GetRegion()
	// If the request has no body with content-length set to 0,
	// we do not have to validate location constraint. Bucket will
	// be created at default region.
	if r.ContentLength == 0 {
		return ErrNone
	}
	locationConstraint := createBucketLocationConfiguration{}
	if err := xmlDecoder(r.Body, &locationConstraint, r.ContentLength); err != nil {
		if err == io.EOF && r.ContentLength == -1 {
			// EOF is a valid condition here when ContentLength is -1.
			return ErrNone
		}
		errorIf(err, "Unable to xml decode location constraint")
		// Treat all other failures as XML parsing errors.
		return ErrMalformedXML
	} // Successfully decoded, proceed to verify the region.

	// Once region has been obtained we proceed to verify it.
	incomingRegion := locationConstraint.Location
	if incomingRegion == "" {
		// Location constraint is empty for region "us-east-1",
		// in accordance with protocol.
		incomingRegion = "us-east-1"
	}
	// Return errInvalidRegion if location constraint does not match
	// with configured region.
	s3Error = ErrNone
	if serverRegion != incomingRegion {
		s3Error = ErrInvalidRegion
	}
	return s3Error
}
