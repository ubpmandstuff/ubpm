// Package gui provides the graphical user interface for ubpm
package gui

import (
	"dura5ka/ubpm/internal/vault"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container" // you'll need that @Kikuoku
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func MakeWindow() {
	myApp := app.New()
	myWindow := myApp.NewWindow("ubpm - usb-based password manager")
	myWindow.Resize(fyne.NewSize(800, 600))
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Vault file path")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Master Password")

	loginButton := widget.NewButton("Open and Unlock Vault", func() {
		v, err := vault.Open(pathEntry.Text, []byte(passwordEntry.Text))
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		showMainManagerScreen(myWindow, v)
	})

	loginScreen := container.NewVBox(
		widget.NewLabel("Please provide vault details to unlock:"),
		pathEntry,
		passwordEntry,
		loginButton,
	)

	myWindow.SetContent(loginScreen)
	myWindow.ShowAndRun()
}

func showMainManagerScreen(myWindow fyne.Window, v *vault.Vault) {
	entries := v.Entries()
	var selectedIndex int = -1

	entryList := widget.NewList(
		func() int {
			return len(entries)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(entries[i].Title)
		},
	)
	entryList.OnSelected = func(id widget.ListItemID) {
		selectedIndex = id
		entry := entries[id]
		dialog.ShowInformation(
			entry.Title,
			fmt.Sprintf("Username: %s\nPassword: %s\nNotes: %s", entry.Username, entry.Password, entry.Notes),
			myWindow,
		)
	}

	addButton := widget.NewButton("Add New Entry", func() {
		titleEntry := widget.NewEntry()
		userEntry := widget.NewEntry()
		passEntry := widget.NewPasswordEntry()
		notesEntry := widget.NewMultiLineEntry()
		formItem := []*widget.FormItem{
			{Text: "Title", Widget: titleEntry},
			{Text: "Username", Widget: userEntry},
			{Text: "Password", Widget: passEntry},
			{Text: "Notes", Widget: notesEntry},
		}

		dialog.ShowForm("Add New Password", "Add", "Cancel", formItem, func(confirm bool) {
			if !confirm {
				return
			}
			err := v.AddEntry(
				titleEntry.Text,
				userEntry.Text,
				passEntry.Text,
				notesEntry.Text,
			)
			if err != nil {
				dialog.ShowError(err, myWindow)
			} else {
				entries = v.Entries()
				entryList.Refresh()
			}
		}, myWindow)
	})

	removeButton := widget.NewButton("Remove Entry", func() {
		if selectedIndex < 0 || selectedIndex >= len(entries) {
			dialog.ShowInformation("No Selection", "Please select an entry to remove.", myWindow)
			return
		}
		entryID := entries[selectedIndex].ID
		err := v.RemoveEntry(entryID)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		entries = v.Entries()
		entryList.Refresh()
		selectedIndex = -1
	})

	buttons := container.NewHBox(addButton, removeButton)

	mainContent := container.NewBorder(
		container.NewHBox(widget.NewLabel("Password Entries:")),
		buttons,
		nil,
		nil,
		entryList,
	)

	myWindow.SetContent(mainContent)
	myWindow.Content().Refresh()
}
