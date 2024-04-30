package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// will be populated by the build command
var version string
var binary_name string

func main() {
	if os.Getenv("CURR_DEV_ENV") == "dev" {
		version = "v0.1.0"
		binary_name = "eddy"
	}

	app := tview.NewApplication()

	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	tree.SetBorder(true)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		// TODO: add a check for files
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(true)
				// SetSelectable(file.IsDir())
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			} else {
				node.SetSelectedFunc(func() {
					node.SetText(node.GetText() + " selected")
				})
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	textArea := tview.NewTextArea()
	textArea.SetTitle(fmt.Sprintf("%s@%s", binary_name, version)).SetBorder(true)

	helpInfo := tview.NewTextView().SetText("Press F1 for help, press Ctrl-C to exit")

	position := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)
	pages := tview.NewPages()

	textArea.SetMovedFunc(func() { updateCursorInfo(position, textArea) })
	updateCursorInfo(position, textArea)

	mainView := tview.NewGrid().
		SetRows(0, 1).
		AddItem(tree, 0, 0, 1, 1, 0, 0, false).
		AddItem(textArea, 0, 1, 1, 2, 0, 0, true).
		AddItem(helpInfo, 1, 0, 1, 1, 0, 0, false).
		AddItem(position, 1, 2, 1, 1, 0, 0, false)

	help1 := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[green]Navigation

[yellow]Left arrow[white]: Move left.
[yellow]Right arrow[white]: Move right.
[yellow]Down arrow[white]: Move down.
[yellow]Up arrow[white]: Move up.
[yellow]Ctrl-A, Home[white]: Move to the beginning of the current line.
[yellow]Ctrl-E, End[white]: Move to the end of the current line.
[yellow]Ctrl-F, page down[white]: Move down by one page.
[yellow]Ctrl-B, page up[white]: Move up by one page.
[yellow]Alt-Up arrow[white]: Scroll the page up.
[yellow]Alt-Down arrow[white]: Scroll the page down.
[yellow]Alt-Left arrow[white]: Scroll the page to the left.
[yellow]Alt-Right arrow[white]:  Scroll the page to the right.
[yellow]Alt-B, Ctrl-Left arrow[white]: Move back by one word.
[yellow]Alt-F, Ctrl-Right arrow[white]: Move forward by one word.

[blue]Press Enter for more help, press Escape to return.`,
		)

	help2 := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[green]Editing[white]

Type to enter text.
[yellow]Ctrl-H, Backspace[white]: Delete the left character.
[yellow]Ctrl-D, Delete[white]: Delete the right character.
[yellow]Ctrl-K[white]: Delete until the end of the line.
[yellow]Ctrl-W[white]: Delete the rest of the word.
[yellow]Ctrl-U[white]: Delete the current line.

[blue]Press Enter for more help, press Escape to return.`,
		)

	help3 := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[green]Selecting Text[white]

Move while holding Shift or drag the mouse.
Double-click to select a word.

[green]Clipboard

[yellow]Ctrl-Q[white]: Copy.
[yellow]Ctrl-X[white]: Cut.
[yellow]Ctrl-V[white]: Paste.
		
[green]Undo

[yellow]Ctrl-Z[white]: Undo.
[yellow]Ctrl-Y[white]: Redo.

[blue]Press Enter for more help, press Escape to return.`,
		)

	help := tview.NewFrame(help1).SetBorders(1, 1, 0, 0, 2, 2)

	help.SetBorder(true).
		SetTitle("Help").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				pages.SwitchToPage("main")
				return nil
			} else if event.Key() == tcell.KeyEnter {
				switch { // can use an array and the modulo operator
				case help.GetPrimitive() == help1:
					help.SetPrimitive(help2)
				case help.GetPrimitive() == help2:
					help.SetPrimitive(help3)
				case help.GetPrimitive() == help3:
					help.SetPrimitive(help1)
				}
				return nil
			}
			return event
		})

	pages.AddAndSwitchToPage("main", mainView, true).
		AddPage("help", tview.NewGrid().
			SetColumns(0, 64, 0).
			SetRows(0, 22, 0).
			AddItem(help, 1, 1, 1, 1, 0, 0, true), true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			pages.ShowPage("help") //TODO: Check when clicking outside help window with the mouse. Then clicking help again.
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages,
		true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
func updateCursorInfo(position *tview.TextView, textArea *tview.TextArea) {
	fromRow, fromColumn, toRow, toColumn := textArea.GetCursor()
	if fromRow == toRow && fromColumn == toColumn {
		position.SetText(fmt.Sprintf("Row: [yellow]%d[white], Column: [yellow]%d ", fromRow, fromColumn))
	} else {
		position.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
	}
}
