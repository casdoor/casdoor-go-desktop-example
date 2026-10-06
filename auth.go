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

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"golang.org/x/oauth2"
)

var casdoorClient = casdoorsdk.NewClient(casdoorEndpoint, clientId, "", certificate, organizationName, applicationName)

// no client secret: a desktop app can't keep one, the authorization code is protected by PKCE instead
var oauthConfig = oauth2.Config{
	ClientID: clientId,
	Endpoint: oauth2.Endpoint{
		AuthURL:   casdoorEndpoint + "/login/oauth/authorize",
		TokenURL:  casdoorEndpoint + "/api/login/oauth/access_token",
		AuthStyle: oauth2.AuthStyleInParams,
	},
	RedirectURL: redirectUri,
	Scopes:      []string{"profile"},
}

func randomString() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// signin opens the Casdoor sign-in page in the browser, waits for Casdoor to redirect
// back to the redirect URI on this machine, and returns the user of the access token.
func signin(ctx context.Context, openBrowser func(u *url.URL) error) (*casdoorsdk.Claims, error) {
	redirectUrl, err := url.Parse(redirectUri)
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", redirectUrl.Host)
	if err != nil {
		return nil, err
	}

	state, err := randomString()
	if err != nil {
		listener.Close()
		return nil, err
	}
	verifier := oauth2.GenerateVerifier()

	type result struct {
		code string
		err  error
	}
	results := make(chan result, 1)

	mux := http.NewServeMux()
	mux.HandleFunc(redirectUrl.Path, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		res := result{code: query.Get("code")}
		if query.Get("state") != state {
			res.err = errors.New("invalid state, please sign in again")
		} else if query.Get("error") != "" {
			res.err = errors.New(strings.TrimSpace(query.Get("error") + " " + query.Get("error_description")))
		} else if res.code == "" {
			res.err = errors.New("no code in the callback")
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if res.err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "<h3>Failed to sign in, you can close this page and try again.</h3>")
		} else {
			fmt.Fprint(w, "<h3>Signed in, you can close this page and go back to the app.</h3>")
		}

		select {
		case results <- res:
		default:
		}
	})
	server := &http.Server{Handler: mux}
	go server.Serve(listener)
	defer server.Close()

	authUrl, err := url.Parse(oauthConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)))
	if err != nil {
		return nil, err
	}
	err = openBrowser(authUrl)
	if err != nil {
		return nil, err
	}

	var res result
	select {
	case res = <-results:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if res.err != nil {
		return nil, res.err
	}

	token, err := oauthConfig.Exchange(ctx, res.code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(token.AccessToken, "error:") {
		return nil, errors.New(strings.TrimSpace(strings.TrimPrefix(token.AccessToken, "error:")))
	}

	// the access token is a JWT, verified with the certificate
	return casdoorClient.ParseJwtToken(token.AccessToken)
}
