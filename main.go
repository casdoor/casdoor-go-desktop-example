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
	"errors"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

func main() {
	a := app.NewWithID("org.casdoor.example.desktop")
	w := a.NewWindow("Casdoor Go Desktop Example")
	w.Resize(fyne.NewSize(480, 320))

	title := widget.NewLabelWithStyle("Casdoor Go Desktop Example", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	message := widget.NewLabel("")
	message.Alignment = fyne.TextAlignCenter
	message.Wrapping = fyne.TextWrapWord

	name := widget.NewLabel("")
	displayName := widget.NewLabel("")
	email := widget.NewLabel("")
	organization := widget.NewLabel("")
	userForm := widget.NewForm(
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Display name", displayName),
		widget.NewFormItem("Email", email),
		widget.NewFormItem("Organization", organization),
	)

	var signinButton, cancelButton, signoutButton *widget.Button
	var cancelSignin context.CancelFunc

	showSignedOut := func(text string) {
		message.SetText(text)
		userForm.Hide()
		signoutButton.Hide()
		cancelButton.Hide()
		signinButton.Show()
	}

	showUser := func(claims *casdoorsdk.Claims) {
		name.SetText(claims.Name)
		displayName.SetText(claims.DisplayName)
		email.SetText(claims.Email)
		organization.SetText(claims.Owner)
		message.SetText("")
		cancelButton.Hide()
		signinButton.Hide()
		userForm.Show()
		signoutButton.Show()
	}

	signinButton = widget.NewButton("Login with Casdoor", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		cancelSignin = cancel
		message.SetText("Sign in with Casdoor in the browser...")
		signinButton.Hide()
		cancelButton.Show()

		go func() {
			defer cancel()
			claims, err := signin(ctx, a.OpenURL)
			fyne.Do(func() {
				if errors.Is(err, context.Canceled) {
					showSignedOut("")
				} else if err != nil {
					showSignedOut("Failed to sign in: " + err.Error())
				} else {
					showUser(claims)
					w.RequestFocus()
				}
			})
		}()
	})
	cancelButton = widget.NewButton("Cancel", func() {
		cancelSignin()
	})
	signoutButton = widget.NewButton("Logout", func() {
		showSignedOut("")
	})
	showSignedOut("")

	w.SetContent(container.NewVBox(
		title,
		message,
		userForm,
		layout.NewSpacer(),
		container.NewCenter(container.NewHBox(signinButton, cancelButton, signoutButton)),
		layout.NewSpacer(),
	))
	w.ShowAndRun()
}
