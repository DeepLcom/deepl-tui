package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TranslatePage provides widgets to translate input text and specify
// translation options.
type TranslatePage struct {
	tview.Pages

	layout *tview.Grid

	sourceLangDropDown *tview.DropDown
	targetLangDropDown *tview.DropDown

	formalityDropDown *tview.DropDown
	glossaryButton    *tview.Button
	glossaryDialog    *GlossariesDialog
	glossaryVisible   bool
	glossarySelected  func(string)

	inputTextArea  *tview.TextArea
	outputTextView *tview.TextView
}

func newTranslatePage(ui *UI) *TranslatePage {
	page := &TranslatePage{
		Pages:  *tview.NewPages(),
		layout: tview.NewGrid(),
	}

	page.sourceLangDropDown = tview.NewDropDown()

	page.inputTextArea = tview.NewTextArea().
		SetPlaceholder("Type to translate.")

	page.outputTextView = tview.NewTextView().
		SetChangedFunc(func() {
			ui.Draw()
		})

	page.targetLangDropDown = tview.NewDropDown()
	page.formalityDropDown = tview.NewDropDown()
	page.glossaryButton = tview.NewButton("Glossary").
		SetSelectedFunc(func() {
			page.setGlossariesDialogVisibility(!page.glossaryVisible)
		})
	container := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(page.targetLangDropDown, 0, 1, true).
		AddItem(page.formalityDropDown, 14, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(page.glossaryButton, 14, 0, false).
		AddItem(nil, 1, 0, false)

	page.layout.
		SetRows(1, 0).
		SetColumns(0, 0).
		SetBorders(true).
		AddItem(page.sourceLangDropDown, 0, 0, 1, 1, 0, 0, false).
		AddItem(container, 0, 1, 1, 1, 0, 0, false).
		AddItem(page.inputTextArea, 1, 0, 1, 1, 0, 0, true).
		AddItem(page.outputTextView, 1, 1, 1, 1, 0, 0, false)
	page.layout.SetBorderPadding(0, 0, 0, 0)

	page.glossaryDialog = newGlossariesDialog().
		SetAcceptedFunc(func(id string, name string) {
			if page.glossarySelected != nil {
				page.glossarySelected(id)
			}
			page.setGlossariesDialogVisibility(false)
		}).
		SetCancelFunc(func() {
			page.setGlossariesDialogVisibility(false)
		})
	page.glossaryDialog.
		SetTitle("Glossary").
		SetBorder(true)

	page.Pages.AddPage("main", page.layout, true, true)
	page.Pages.AddPage("dialog", page.glossaryDialog, false, true)
	page.Pages.HidePage("dialog")

	page.registerKeyBindings(ui)

	return page
}

// SetGlossarySelectedFunc sets a handler that is called when a glossary is
// selected.
func (w *TranslatePage) SetGlossarySelectedFunc(selected func(string)) *TranslatePage {
	w.glossarySelected = selected
	return w
}

func (w *TranslatePage) setGlossariesDialogVisibility(visible bool) {
	if visible {
		w.Pages.ShowPage("dialog")
	} else {
		w.Pages.HidePage("dialog")
	}
	w.glossaryVisible = visible
}

func (w *TranslatePage) registerKeyBindings(ui *UI) {
	w.layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModAlt {
			switch event.Rune() {
			case 's':
				ui.SetFocus(w.sourceLangDropDown)
				return nil
			case 't':
				ui.SetFocus(w.targetLangDropDown)
				return nil
			case 'f':
				ui.SetFocus(w.formalityDropDown)
				return nil
			case 'g':
				ui.SetFocus(w.glossaryButton)
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

	if width > 156 {
		sourceLangLabel = "Select source language: "
		targetLangLabel = "Select target language: "
	} else {
		sourceLangLabel = ""
		targetLangLabel = ""
	}

	w.sourceLangDropDown.SetLabel(sourceLangLabel)
	w.targetLangDropDown.SetLabel(targetLangLabel)

	var (
		gbx, gby, gbw, _ = w.glossaryButton.GetRect()
		gww, gwh         = 40, 20
	)

	w.glossaryDialog.SetRect(gbx+gbw-gww, gby, gww, gwh)
}
