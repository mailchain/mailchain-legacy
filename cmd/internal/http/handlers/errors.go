// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

// ValidationError describes a 422 validation error.
// swagger:response ValidationError
type ValidationError struct {
	// Code describing the error
	// example: 422
	Code string `json:"code"`

	// Description of the error
	// example: Response to invalid input
	Message string `json:"message"`
}

// NotFoundError describes a 404 not found error.
// swagger:response NotFoundError
type NotFoundError struct {
	// Code describing the error
	// example: 404
	Code string `json:"code"`

	// Description of the error
	// example: Not found.
	Message string `json:"message"`
}
