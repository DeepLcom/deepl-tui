package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	headerTextLarge = `
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
	headerTextMedium = `
░░░░░░░░░░░░▒▓██▓▒░░░░░░░░░░░░
░░░░░░░░▒▒▓████████▓▒▒░░░░░░░░
░░░░░░▓█████▓▓█████████▓░░░░░░  _______                  _       _             
░░░░░░█████▓░░▓█████████░░░░░░ |__   __|                | |     | |            
░░░░░░██████▓█▓▓▓▓▒▓████░░░░░░    | |_ __ __ _ _ __  ___| | __ _| |_ ___  _ __ 
░░░░░░██████████▓▒░▒████░░░░░░    | | '__/ _\ | '_ \/ __| |/ _\ | __/ _ \| '__|
░░░░░░█████▓░▒▒▓████████░░░░░░    | | | | (_| | | | \__ \ | (_| | || (_) | |   
░░░░░░▓████▓▒▒█████████▓░░░░░░    |_|_|  \__,_|_| |_|___/_|\__,_|\__\___/|_|   
░░░░░░░░▒▒▓████████▓▒▒░░░░░░░░
░░░░░░░░░░░░▒▓████▒░░░░░░░░░░░
░░░░░░░░░░░░░░░▒▓█░░░░░░░░░░░░
░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
`
	headerTextSmall = `
░░░░░░░░░░░░▒▓██▓▒░░░░░░░░░░░░
░░░░░░░░▒▒▓████████▓▒▒░░░░░░░░
░░░░░░▓█████▓▓█████████▓░░░░░░
░░░░░░█████▓░░▓█████████░░░░░░
░░░░░░██████▓█▓▓▓▓▒▓████░░░░░░
░░░░░░██████████▓▒░▒████░░░░░░
░░░░░░█████▓░▒▒▓████████░░░░░░
░░░░░░▓████▓▒▒█████████▓░░░░░░
░░░░░░░░▒▒▓████████▓▒▒░░░░░░░░
░░░░░░░░░░░░▒▓████▒░░░░░░░░░░░
░░░░░░░░░░░░░░░▒▓█░░░░░░░░░░░░
░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
`
)

var debug bool = false

type ui struct {
	tview.Application

	grid *tview.Grid

	header *tview.TextView
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

	ui.header = tview.NewTextView().
		SetTextAlign(tview.AlignLeft)

	ui.sourceLangDropDown = tview.NewDropDown()
	ui.targetLangDropDown = tview.NewDropDown()

	ui.inputTextArea = tview.NewTextArea().
		SetPlaceholder("Type to translate.")

	ui.outputTextView = tview.NewTextView()
	ui.outputTextView.SetChangedFunc(func() {
		ui.Draw()
	})

	ui.footer = tview.NewTextView().
		SetTextAlign(tview.AlignRight)

	grid := tview.NewGrid()

	grid.
		AddItem(ui.header, 0, 0, 1, 2, 0, 0, false).
		AddItem(ui.sourceLangDropDown, 1, 0, 1, 1, 0, 0, false).
		AddItem(ui.targetLangDropDown, 1, 1, 1, 1, 0, 0, false).
		AddItem(ui.inputTextArea, 2, 0, 1, 1, 0, 0, true).
		AddItem(ui.outputTextView, 2, 1, 1, 1, 0, 0, false).
		AddItem(ui.footer, 3, 0, 1, 2, 0, 0, false)

	grid.
		SetColumns(0, 0).
		SetBorders(true)

	ui.grid = grid

	ui.SetRoot(ui.grid, true)

	ui.registerKeybindings()

	ui.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		w, h := screen.Size()
		return ui.adjustToScreenSize(w, h)
	})

	return ui
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

func (ui *ui) adjustToScreenSize(width int, height int) bool {
	var (
		headerText   string
		headerHeight int

		sourceLangLabel string
		targetLangLabel string
	)

	if width > 112 && height > 30 {
		headerText = headerTextLarge
		headerHeight = 12
	} else if width > 80 && height > 30 {
		headerText = headerTextMedium
		headerHeight = 12
	} else if width > 31 && height > 30 {
		headerText = headerTextSmall
		headerHeight = 12
	} else {
		headerText = "DeepL Translator"
		headerHeight = 1
	}

	ui.grid.SetRows(headerHeight, 1, 0, 1)
	ui.header.SetText(strings.TrimPrefix(headerText, "\n"))

	if width > 96 {
		sourceLangLabel = "Select source language: "
		targetLangLabel = "Select target language: "
	} else {
		sourceLangLabel = ""
		targetLangLabel = ""
	}

	ui.sourceLangDropDown.SetLabel(sourceLangLabel)
	ui.targetLangDropDown.SetLabel(targetLangLabel)

	if debug {
		ui.footer.SetText(fmt.Sprintf("(%d,%d)", width, height))
	}

	return false
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
