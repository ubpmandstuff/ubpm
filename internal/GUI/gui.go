package gui

import (
	"dura5ka/ubpm/internal/vault"

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
		showMainManagerScreen(myWindow)
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

func showMainManagerScreen(myWindow fyne.Window) {

	listLabel := widget.NewLabel("Welcome! Entries will appear here once loaded.")

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
	})
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
		}
	})

	mainContent := container.NewBorder(
		container.NewHBox(widget.NewLabel("Password Entries:")),
		addButton,
		nil,
		nil,
		listLabel, // Заменить на widget.NewList
	)

	myWindow.SetContent(mainContent)
	myWindow.Content().Refresh()
}
