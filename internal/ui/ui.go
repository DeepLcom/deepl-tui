package ui

import (
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

type UI struct {
	tview.Application

	layout *tview.Grid

	header *tview.TextView
	footer *tview.InputField

	pages          *tview.Pages
	pageIndex      []string
	translatePage  *TranslatePage
	glossariesPage *GlossariesPage
}

func NewUI() *UI {
	ui := &UI{
		Application: *tview.NewApplication(),
	}

	ui.header = tview.NewTextView().
		SetTextAlign(tview.AlignLeft)
	ui.header.SetBorder(true)

	ui.setupFooter()

	ui.translatePage = newTranslatePage(ui)
	ui.glossariesPage = newGlossariesPage(ui)

	ui.pages = tview.NewPages()
	ui.pages.AddPage("translate", ui.translatePage, true, true)
	ui.pageIndex = append(ui.pageIndex, "translate")
	ui.pages.AddPage("glossaries", ui.glossariesPage, true, false)
	ui.pageIndex = append(ui.pageIndex, "glossaries")

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

func (ui *UI) registerKeybindings() {
	ui.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyCtrlC:
            // don't quit here
            return nil
        case tcell.KeyCtrlQ:
            ui.Application.Stop()
        }

		if event.Modifiers() == tcell.ModAlt {
			if event.Key() == tcell.KeyTab {
				ui.cycePage()
				return nil
			}

			switch event.Rune() {
			case ':':
				ui.switchToCommandPrompt()
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

	ui.translatePage.adjustToSize()

	return false
}

func (ui *UI) switchToPage(name string) {
	ui.pages.SwitchToPage(name)
}

func (ui *UI) cycePage() {
	name, _ := ui.pages.GetFrontPage()
	var j int = 0
	for i, n := range ui.pageIndex {
		if n == name {
			if i < ui.pages.GetPageCount()-1 {
				j = i + 1
			}
		}
	}
	ui.switchToPage(ui.pageIndex[j])
}

func (ui *UI) switchToCommandPrompt() {
	ui.footer.
		SetLabel(":").
		SetText("").
		SetDisabled(false)
	ui.SetFocus(ui.footer)
}

func (ui *UI) SetFooter(text string) {
	ui.footer.SetText(text)
}

func (ui *UI) GetCurrentSourceLang() string {
	_, opt := ui.translatePage.sourceLangDropDown.GetCurrentOption()
	return opt
}

func (ui *UI) GetCurrentTargetLang() string {
	_, opt := ui.translatePage.targetLangDropDown.GetCurrentOption()
	return opt
}

func (ui *UI) SetSourceLangOptions(opts []string, selected func(string, int)) {
	ui.translatePage.sourceLangDropDown.
		SetOptions(opts, selected).
		SetCurrentOption(0)
}

func (ui *UI) SetTargetLangOptions(opts []string, selected func(string, int)) {
	ui.translatePage.targetLangDropDown.SetOptions(opts, selected)
}

func (ui *UI) SetInputTextChangedFunc(handler func()) {
	ui.translatePage.inputTextArea.SetChangedFunc(handler)
}

func (ui *UI) GetInputText() string {
	return ui.translatePage.inputTextArea.GetText()
}

func (ui *UI) WriteOutputText(r io.Reader) error {
	_, err := io.Copy(ui.translatePage.outputTextView, r)
	if err != nil {
		return err
	}
	return nil
}

func (ui *UI) ClearOutputText() {
	ui.translatePage.outputTextView.Clear()
}
