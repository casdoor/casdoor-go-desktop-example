// Copyright 2026 The Casdoor Authors. All Rights Reserved.
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

package main

import _ "embed"

// The Casdoor application to sign in with, the defaults are the public demo server https://door.casdoor.com
const (
	casdoorEndpoint  = "https://door.casdoor.com"
	clientId         = "014ae4bd048734ca2dea"
	organizationName = "casbin"
	applicationName  = "app-casnode"
	// must be in the Redirect URLs of the application, the app listens on it while signing in
	redirectUri = "http://localhost:9000/callback"
)

// the certificate of the cert used by the application, see Casdoor -> Certs
//
//go:embed cert.pem
var certificate string
