package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create the left pane: a list for the file explorer.
	explorer := tview.NewList()
	explorer.SetBorder(true)
	explorer.SetTitle("File Explorer")

	// Create the right pane: a list for selected files.
	selectedFiles := tview.NewList()
	selectedFiles.SetBorder(true)
	selectedFiles.SetTitle("Selected Files")

	// Use a Flex container to arrange the two panes side by side.
	flex := tview.NewFlex().
		AddItem(explorer, 0, 2, true).  // Left pane gets focus initially.
		AddItem(selectedFiles, 0, 1, false)

	// Get the current working directory.
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// loadDir reads the directory and populates the explorer.
	var loadDir func(string)
	loadDir = func(dir string) {
		explorer.Clear()

		// Add an entry to go to the parent directory if not at root.
		if dir != "/" {
			explorer.AddItem("..", "    -Parent Directory", 0, func() {
				loadDir(filepath.Dir(dir))
			})
		}

		// Read the directory contents.
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			explorer.AddItem("Error reading directory", "", 0, nil)
			return
		}

		// Add each file and directory to the list.
		for _, file := range files {
			// Create a local copy to avoid closure issues.
			f := file
			if f.IsDir() {
				// For directories, append "/" to the name.
				explorer.AddItem(f.Name()+"/", "    -Directory", 0, func() {
					loadDir(filepath.Join(dir, f.Name()))
				})
			} else {
				// For files, pressing Enter will add the file (with its full path) to the right list.
				explorer.AddItem(f.Name(), "    -File", 0, func() {
					fullPath := filepath.Join(dir, f.Name())
					selectedFiles.AddItem(fullPath, "", 0, nil)
				})
			}
		}

		// Update the left pane title to show the current directory.
		explorer.SetTitle("File Explorer - " + dir)
	}

	// Load the initial directory.
	loadDir(currentDir)

	// Allow switching focus between the left and right panes with the Tab key,
	// and open a modal to change directory with Ctrl+D.
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if app.GetFocus() == explorer {
				app.SetFocus(selectedFiles)
			} else {
				app.SetFocus(explorer)
			}
			return nil
		case tcell.KeyCtrlD:
			// Create the form outside of the method chain so that it can be captured in closures.
			directoryForm := tview.NewForm()
			directoryForm.AddInputField("Directory", "", 40, nil, nil)
			directoryForm.AddButton("Change", func() {
				newDir := directoryForm.GetFormItemByLabel("Directory").(*tview.InputField).GetText()
				loadDir(newDir)
				app.SetRoot(flex, true)
				app.SetFocus(explorer)
			})
			directoryForm.AddButton("Cancel", func() {
				app.SetRoot(flex, true)
				app.SetFocus(explorer)
			})
			directoryForm.SetBorder(true)
			directoryForm.SetTitle("Change Directory")
			directoryForm.SetTitleAlign(tview.AlignLeft)
			app.SetRoot(directoryForm, true).SetFocus(directoryForm)
			return nil
		}
		return event
	})

	// Set the Flex container as the root and run the application.
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		log.Fatal(err)
	}
}
