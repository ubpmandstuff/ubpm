package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container" // you'll need that @Kikuoku
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
		fmt.Println("Placeholder for unlocking logic.") // Заглушка
		showMainManagerScreen(myWindow)                 // Заглушка
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

func showMainManagerScreen(myWindow fyne.Window) { // Заглушка

	listLabel := widget.NewLabel("Welcome! Entries will appear here once loaded.")

	addButton := widget.NewButton("Add New Entry", func() {
		fmt.Println("Add entry UI logic needed.")
	})

	// saveButton := widget.NewButton("Save All Changes", func() {
	// 	fmt.Println("Save logic needed.")
	// })

	mainContent := container.NewBorder(
		// container.NewHBox(widget.NewLabel("Password Entries:"), layout.NewSpacer(), saveButton),
		container.NewHBox(widget.NewLabel("Password Entries:")),
		addButton,
		nil,
		nil,
		listLabel, // Заменить на widget.NewList
	)

	myWindow.SetContent(mainContent)
	myWindow.Content().Refresh()
}
