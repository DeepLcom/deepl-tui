package ui

import (
	"fmt"
	"io"
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

var debug bool = true

type UI struct {
	tview.Application

	layout *tview.Grid

	header *tview.TextView
	footer *tview.TextView
	pages  *tview.Pages

	// translate page
	sourceLangDropDown *tview.DropDown
	targetLangDropDown *tview.DropDown

	inputTextArea  *tview.TextArea
	outputTextView *tview.TextView

	// glossaries page

}

func NewUI() *UI {
	ui := &UI{
		Application: *tview.NewApplication(),
	}

	ui.header = tview.NewTextView().
		SetTextAlign(tview.AlignLeft)
	ui.header.SetBorder(true)

	ui.footer = tview.NewTextView().
		SetTextAlign(tview.AlignRight)
	ui.footer.SetBorder(true)

	ui.pages = tview.NewPages().
		AddPage("translate", ui.setupTranslatePage(), true, true).
		AddPage("glossaries", ui.setupGlossariesPage(), true, false)

	ui.layout = tview.NewGrid().
		SetBorders(false).
		AddItem(ui.header, 0, 0, 1, 1, 0, 0, false).
		AddItem(ui.pages, 1, 0, 1, 1, 0, 0, true).
		AddItem(ui.footer, 2, 0, 1, 1, 0, 0, false)

	ui.SetRoot(ui.layout, true)

	ui.registerKeybindings()

	ui.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		w, h := screen.Size()
		return ui.adjustToScreenSize(w, h)
	})

	return ui
}

func (ui *UI) setupTranslatePage() tview.Primitive {
	ui.sourceLangDropDown = tview.NewDropDown()
	ui.targetLangDropDown = tview.NewDropDown()

	ui.inputTextArea = tview.NewTextArea().
		SetPlaceholder("Type to translate.")

	ui.outputTextView = tview.NewTextView()
	ui.outputTextView.SetChangedFunc(func() {
		ui.Draw()
	})

	layout := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(0, 0).
		SetBorders(true).
		AddItem(ui.sourceLangDropDown, 0, 0, 1, 1, 0, 0, false).
		AddItem(ui.targetLangDropDown, 0, 1, 1, 1, 0, 0, false).
		AddItem(ui.inputTextArea, 1, 0, 1, 1, 0, 0, true).
		AddItem(ui.outputTextView, 1, 1, 1, 1, 0, 0, false)
	layout.SetBorderPadding(0, 0, 0, 0)

	return layout
}

func (ui *UI) setupGlossariesPage() tview.Primitive {
	layout := tview.NewGrid()
	return layout
}

func (ui *UI) registerKeybindings() {
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

func (ui *UI) adjustToScreenSize(width int, height int) bool {
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

	ui.layout.SetRows(headerHeight+2, 0, 1+2)
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

func (ui *UI) SetFooter(text string) {
	ui.footer.SetText(text)
}

func (ui *UI) GetCurrentSourceLang() string {
	_, opt := ui.sourceLangDropDown.GetCurrentOption()
	return opt
}

func (ui *UI) GetCurrentTargetLang() string {
	_, opt := ui.targetLangDropDown.GetCurrentOption()
	return opt
}

func (ui *UI) SetSourceLangOptions(opts []string, selected func(string, int)) {
	ui.sourceLangDropDown.
		SetOptions(opts, selected).
		SetCurrentOption(0)
}

func (ui *UI) SetTargetLangOptions(opts []string, selected func(string, int)) {
	ui.targetLangDropDown.SetOptions(opts, selected)
}

func (ui *UI) SetInputTextChangedFunc(handler func()) {
	ui.inputTextArea.SetChangedFunc(handler)
}

func (ui *UI) GetInputText() string {
	return ui.inputTextArea.GetText()
}

func (ui *UI) WriteOutputText(r io.Reader) error {
	_, err := io.Copy(ui.outputTextView, r)
	if err != nil {
		return err
	}
	return nil
}

func (ui *UI) ClearOutputText() {
	ui.outputTextView.Clear()
}
