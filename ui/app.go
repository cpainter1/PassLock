package ui

import (
	"fyne.io/fyne/v2/app"
	"github.com/cpainter1/PassLock/ui/views"
)

func RunApp() {
	myApp := app.New()
	mainWindow := myApp.NewWindow("PassLock")

	ui.ShowLoginUI(mainWindow)

	mainWindow.ShowAndRun()
}
