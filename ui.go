package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ui struct {
	tview.Application

	header tview.Primitive
	footer *tview.TextView

	inputTextArea  *tview.TextArea
	outputTextView *tview.TextView

	sourceLangDropDown *tview.DropDown
	targetLangDropDown *tview.DropDown
}

func newUI() *ui {
	ui := &ui{
		Application: *tview.NewApplication(),
	}

	ui.setupHeader()

	ui.sourceLangDropDown = tview.NewDropDown().
		SetLabel("Select source language: ")

	ui.targetLangDropDown = tview.NewDropDown().
		SetLabel("Select target language: ")

	ui.inputTextArea = tview.NewTextArea().
		SetPlaceholder("Type to translate.")

	ui.outputTextView = tview.NewTextView()
	ui.outputTextView.SetChangedFunc(func() {
		ui.Draw()
	})

	ui.footer = tview.NewTextView()

	grid := tview.NewGrid()

	grid.
		AddItem(ui.header, 0, 0, 1, 2, 0, 0, false).
		AddItem(ui.sourceLangDropDown, 1, 0, 1, 1, 0, 0, false).
		AddItem(ui.targetLangDropDown, 1, 1, 1, 1, 0, 0, false).
		AddItem(ui.inputTextArea, 2, 0, 1, 1, 0, 0, true).
		AddItem(ui.outputTextView, 2, 1, 1, 1, 0, 0, false).
		AddItem(ui.footer, 3, 0, 1, 2, 0, 0, false)

	grid.
		SetRows(12, 1, 0, 1).
		SetColumns(0, 0).
		SetBorders(true)

	ui.SetRoot(grid, true)

	ui.registerKeybindings()

	return ui
}

func (ui *ui) setupHeader() {
	var text string

	text = `
░░░░░░░░░░░░▒▓██▓▒░░░░░░░░░░░░
░░░░░░░░▒▒▓████████▓▒▒░░░░░░░░
░░░░░░▓█████▓▓█████████▓░░░░░░  _____                  _        _______                  _       _             
░░░░░░█████▓░░▓█████████░░░░░░ |  __ \                | |      |__   __|                | |     | |            
░░░░░░██████▓█▓▓▓▓▒▓████░░░░░░ | |  | | ___  ___ _ __ | |         | |_ __ __ _ _ __  ___| | __ _| |_ ___  _ __ 
░░░░░░██████████▓▒░▒████░░░░░░ | |  | |/ _ \/ _ \ '_ \| |         | | '__/ _\ | '_ \/ __| |/ _\ | __/ _ \| '__|
░░░░░░█████▓░▒▒▓████████░░░░░░ | |__| |  __/  __/ |_) | |____     | | | | (_| | | | \__ \ | (_| | || (_) | |   
░░░░░░▓████▓▒▒█████████▓░░░░░░ |_____/ \___|\___| .__/|______|    |_|_|  \__,_|_| |_|___/_|\__,_|\__\___/|_|   
░░░░░░░░▒▒▓████████▓▒▒░░░░░░░░                  | |                                                            
░░░░░░░░░░░░▒▓████▒░░░░░░░░░░░                  |_|
░░░░░░░░░░░░░░░▒▓█░░░░░░░░░░░░
░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
`

	ui.header = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(strings.TrimSpace(text))
}

func (ui *ui) SetFooter(text string) {
	ui.footer.SetText(text)
}

func (ui *ui) GetCurrentSourceLang() string {
	_, opt := ui.sourceLangDropDown.GetCurrentOption()
	return opt
}

func (ui *ui) GetCurrentTargetLang() string {
	_, opt := ui.targetLangDropDown.GetCurrentOption()
	return opt
}

func (ui *ui) SetSourceLangOptions(opts []string, selected func(string, int)) {
	ui.sourceLangDropDown.
		SetOptions(opts, selected).
		SetCurrentOption(0)
}

func (ui *ui) SetTargetLangOptions(opts []string, selected func(string, int)) {
	ui.targetLangDropDown.SetOptions(opts, selected)
}
func (ui *ui) registerKeybindings() {
	ui.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModAlt {
			switch event.Rune() {
			case 's':
				ui.SetFocus(ui.sourceLangDropDown)
				return nil
			case 't':
				ui.SetFocus(ui.targetLangDropDown)
				return nil
			case 'i':
				ui.SetFocus(ui.inputTextArea)
				return nil
			}
		}
		return event
	})
}
