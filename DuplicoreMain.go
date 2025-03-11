package main

import (
	"fmt"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create the form and store it in a variable.
	form := tview.NewForm()
	form.AddInputField("Script Name", "", 20, nil, nil)
	form.AddInputField("Install Command", "", 40, nil, nil)
	form.AddButton("Generate", func() {
		// Retrieve the text from the input fields.
		scriptName := form.GetFormItemByLabel("Script Name").(*tview.InputField).GetText()
		command := form.GetFormItemByLabel("Install Command").(*tview.InputField).GetText()

		// Generate a simple bash install script.
		script := fmt.Sprintf("#!/bin/bash\n\n# %s\n%s\n", scriptName, command)

		// Display the generated script in a modal window.
		modal := tview.NewModal()
		modal.SetText(script).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.SetRoot(form, true)
			})
		app.SetRoot(modal, true)
	})
	form.AddButton("Quit", func() {
		app.Stop()
	})

	form.SetBorder(true)
	form.SetTitle("Install Script Creator")
	form.SetTitleAlign(tview.AlignLeft)

	// Start the application.
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
