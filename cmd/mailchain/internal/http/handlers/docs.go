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

// Package handlers contains the API handers for the mailchain client.
package handlers

import (
	"net/http"
)

// GetDocs returns a handler get spec
func GetDocs() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		  <head>
			<title>ReDoc</title>
			<!-- needed for adaptive design -->
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1">
		
			<!--
			ReDoc doesn't change outer page styles
			-->
			<style>
			  body {
				margin: 0;
				padding: 0;
			  }
			</style>
		  </head>
		  <body>
			<redoc spec-url='./spec.json'></redoc>
			<script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
		  </body>
		</html>
		`))
	}
}
