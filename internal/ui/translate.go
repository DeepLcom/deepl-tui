package ui

import (
	"github.com/cluttrdev/deepl-go/deepl"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TranslatePage struct {
	tview.Pages

	layout *tview.Grid

	sourceLangDropDown *tview.DropDown
	targetLangDropDown *tview.DropDown

	formalityDropDown *tview.DropDown
	glossaryButton    *tview.Button
	glossaryWidget    *GlossariesWidget
	glossarySelected  func(string, int)
	glossaryDialog    *Dialog
	glossaryVisible   bool

	inputTextArea  *tview.TextArea
	outputTextView *tview.TextView
}

func newTranslatePage(ui *UI) *TranslatePage {
	widget := &TranslatePage{
		Pages:  *tview.NewPages(),
		layout: tview.NewGrid(),
	}

	widget.sourceLangDropDown = tview.NewDropDown()

	widget.inputTextArea = tview.NewTextArea().
		SetPlaceholder("Type to translate.")

	widget.outputTextView = tview.NewTextView().
		SetChangedFunc(func() {
			ui.Draw()
		})

	widget.targetLangDropDown = tview.NewDropDown()
	widget.formalityDropDown = tview.NewDropDown()
	widget.glossaryButton = tview.NewButton("Glossary").
		SetSelectedFunc(func() {
			widget.setGlossariesDialogVisibility(!widget.glossaryVisible)
		})
	container := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(widget.targetLangDropDown, 0, 1, true).
		AddItem(widget.formalityDropDown, 14, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(widget.glossaryButton, 14, 0, false).
		AddItem(nil, 1, 0, false)

	widget.layout.
		SetRows(1, 0).
		SetColumns(0, 0).
		SetBorders(true).
		AddItem(widget.sourceLangDropDown, 0, 0, 1, 1, 0, 0, false).
		AddItem(container, 0, 1, 1, 1, 0, 0, false).
		AddItem(widget.inputTextArea, 1, 0, 1, 1, 0, 0, true).
		AddItem(widget.outputTextView, 1, 1, 1, 1, 0, 0, false)
	widget.layout.SetBorderPadding(0, 0, 0, 0)

	widget.glossaryWidget = newGlossariesWidget(ui).
		SetSelectedFunc(func(text string, index int) {

		})
	widget.glossaryDialog = NewDialog(widget.glossaryWidget).
		AddButton("OK", func() {
			index, text := widget.glossaryWidget.GetCurrentOption()
			if widget.glossarySelected != nil {
				widget.glossarySelected(text, index)
			}
			widget.setGlossariesDialogVisibility(false)
		}).
		SetCancelFunc(func() {
			widget.setGlossariesDialogVisibility(false)
		})
	widget.glossaryDialog.
		SetTitle("Glossary").
		SetBorder(true)

	widget.Pages.AddPage("main", widget.layout, true, true)
	widget.Pages.AddPage("dialog", widget.glossaryDialog, false, true)
	widget.Pages.HidePage("dialog")

	widget.registerKeyBindings(ui)

	return widget
}

func (w *TranslatePage) SetGlossaryDataFunc(data func(string, int) (*deepl.GlossaryInfo, []deepl.GlossaryEntry)) *TranslatePage {
	w.glossaryWidget.SetDataFunc(data)
	return w
}

func (w *TranslatePage) SetGlossarySelectedFunc(selected func(string, int)) *TranslatePage {
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
			case 'g':
				ui.SetFocus(w.glossaryButton)
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
