package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Button Widget")

	content := widget.NewButton("Hi", func() {
		log.Println("Here we go!")
	})
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
