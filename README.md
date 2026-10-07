# Casdoor Go Desktop Example

[![Build](https://github.com/casdoor/casdoor-go-desktop-example/actions/workflows/build.yml/badge.svg)](https://github.com/casdoor/casdoor-go-desktop-example/actions/workflows/build.yml)
[![License](https://img.shields.io/github/license/casdoor/casdoor-go-desktop-example)](https://github.com/casdoor/casdoor-go-desktop-example/blob/master/LICENSE)
[![Discord](https://img.shields.io/discord/1022748306096537660?logo=discord&label=discord&color=5865F2)](https://discord.gg/5rPsrAzK7S)

An example cross-platform desktop app in Go, built with [Fyne](https://fyne.io/), that signs users in with [Casdoor](https://casdoor.ai/) in the system browser, using the OAuth 2.0 authorization code flow with PKCE and a loopback redirect URI ([RFC 8252](https://datatracker.ietf.org/doc/html/rfc8252)).

## How it works

1. **Login with Casdoor** starts a small HTTP server on the redirect URI `http://localhost:9000/callback`, creates a PKCE code verifier and a random state, and opens the Casdoor sign-in page in the browser ([auth.go](auth.go)).
2. After signing in, Casdoor redirects the browser to `http://localhost:9000/callback?code=...&state=...`, which the app is listening on. The app checks the state.
3. The app exchanges the code for the tokens with the code verifier ([golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)). No client secret is stored in the app.
4. The access token is a JWT: [casdoor-go-sdk](https://github.com/casdoor/casdoor-go-sdk) verifies it with the certificate of the application (`ParseJwtToken()`) and returns the user, which the window shows ([main.go](main.go)).

## Prerequisites

- Go 1.23+
- A C compiler and the graphics libraries for Fyne, see [Fyne: getting started](https://docs.fyne.io/started/). For example on Ubuntu: `sudo apt install gcc libgl1-mesa-dev xorg-dev libwayland-dev libxkbcommon-dev`; on Windows: [MSYS2](https://www.msys2.org/) or [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
- A Casdoor server. The example is preconfigured for the public demo server https://door.casdoor.com, so it runs as is. To use your own, see [Casdoor installation](https://casdoor.ai/docs/basic/server-installation).

## Configuration

Skip this section to try the example with the public demo server.

In your Casdoor, create (or reuse) an organization and an application, and add `http://localhost:9000/callback` to the application's **Redirect URLs**. Then fill in [conf.go](conf.go):

```go
const (
	casdoorEndpoint  = "https://door.casdoor.com" // Casdoor server URL
	clientId         = "014ae4bd048734ca2dea"     // client ID of the application
	organizationName = "casbin"                   // organization of the application
	applicationName  = "app-casnode"              // name of the application
	redirectUri      = "http://localhost:9000/callback"
)
```

and put the certificate of the cert used by the application (Casdoor -> Certs) in [cert.pem](cert.pem).

## Run

```shell
git clone https://github.com/casdoor/casdoor-go-desktop-example
cd casdoor-go-desktop-example
go run .
```

Click **Login with Casdoor**. On the demo server, sign in with username `admin` and password `123`, then go back to the app.

To package the app with an icon and metadata, see [fyne package](https://docs.fyne.io/started/packaging).

## Resources

- [Casdoor documentation](https://casdoor.ai/docs/overview)
- [casdoor-go-sdk](https://github.com/casdoor/casdoor-go-sdk)
- [Fyne](https://fyne.io/)

## License

[Apache-2.0](LICENSE)
