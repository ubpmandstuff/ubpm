package gui

import (
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func MakeWindow() fyne.Window {
	myApp := app.New()
	myWindow := myApp.NewWindow("Button Widget")

	content := widget.NewButton("Hi", func() {
		log.Println("Here we go!")
	})
	myWindow.SetContent(content)
	return myWindow
}
