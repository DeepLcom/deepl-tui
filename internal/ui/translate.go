package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TranslatePage struct {
	tview.Grid

	sourceLangDropDown *tview.DropDown
	targetLangDropDown *tview.DropDown

	inputTextArea  *tview.TextArea
	outputTextView *tview.TextView
}

func newTranslatePage(ui *UI) *TranslatePage {
	page := &TranslatePage{
		Grid: *tview.NewGrid(),
	}

	page.sourceLangDropDown = tview.NewDropDown()
	page.targetLangDropDown = tview.NewDropDown()

	page.inputTextArea = tview.NewTextArea().
		SetPlaceholder("Type to translate.")

	page.outputTextView = tview.NewTextView()
	page.outputTextView.SetChangedFunc(func() {
		ui.Draw()
	})

	page.Grid.
		SetRows(1, 0).
		SetColumns(0, 0).
		SetBorders(true).
		AddItem(page.sourceLangDropDown, 0, 0, 1, 1, 0, 0, false).
		AddItem(page.targetLangDropDown, 0, 1, 1, 1, 0, 0, false).
		AddItem(page.inputTextArea, 1, 0, 1, 1, 0, 0, true).
		AddItem(page.outputTextView, 1, 1, 1, 1, 0, 0, false)
	page.Grid.SetBorderPadding(0, 0, 0, 0)

	page.registerKeyBindings(ui)

	return page
}

func (w *TranslatePage) registerKeyBindings(ui *UI) {
	w.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModAlt {
			switch event.Rune() {
			case 's':
				ui.SetFocus(w.sourceLangDropDown)
				return nil
			case 't':
				ui.SetFocus(w.targetLangDropDown)
				return nil
			case 'i':
				ui.SetFocus(w.inputTextArea)
				return nil
			}
		}
		return event
	})

}

func (w *TranslatePage) adjustToSize() {
	_, _, width, _ := w.GetInnerRect()

	var (
		sourceLangLabel string
		targetLangLabel string
	)

	if width > 96 {
		sourceLangLabel = "Select source language: "
		targetLangLabel = "Select target language: "
	} else {
		sourceLangLabel = ""
		targetLangLabel = ""
	}

	w.sourceLangDropDown.SetLabel(sourceLangLabel)
	w.targetLangDropDown.SetLabel(targetLangLabel)
}
