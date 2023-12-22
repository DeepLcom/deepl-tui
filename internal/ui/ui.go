package ui

import (
	"io"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/cluttrdev/deepl-go/deepl"
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
	_, page := ui.pages.GetFrontPage()
	ui.SetFocus(page)
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

func (ui *UI) SetSourceLangOptions(opts []string, selected func(string, int)) {
	ui.translatePage.sourceLangDropDown.
		SetOptions(opts, selected).
		SetCurrentOption(0)
}

func (ui *UI) SetTargetLangOptions(opts []string, selected func(string, int)) {
	ui.translatePage.targetLangDropDown.SetOptions(opts, selected)
}

func (ui *UI) SetFormalityOptions(opts []string, selected func(string, int)) {
	ui.translatePage.formalityDropDown.
		SetOptions(opts, selected).
		SetCurrentOption(0)
}

// SetGlossaryOptions provides the ui with a list of available glossary ids and names.
func (ui *UI) SetGlossaryOptions(options [][2]string) {
	ui.translatePage.glossaryDialog.SetOptions(options)
	ui.glossariesPage.SetOptions(options)
}

// SetGlossaryLanguageOptions provides the ui with a list of supported glossary languages.
func (ui *UI) SetGlossaryLanguageOptions(langs []string) {
	ui.glossariesPage.SetLanguageOptions(langs)
}

// SetGlossaryDataFunc sets a handler which can be called by the ui widgets to request glossary data.
// It receives the glossary ID as an argument.
func (ui *UI) SetGlossaryDataFunc(handler func(string) (deepl.GlossaryInfo, []deepl.GlossaryEntry)) {
	ui.translatePage.glossaryDialog.SetDataFunc(handler)
	ui.glossariesPage.SetGlossaryDataFunc(handler)
}

func (ui *UI) SetGlossarySelectedFunc(handler func(string)) {
	ui.translatePage.SetGlossarySelectedFunc(handler)
}

func (ui *UI) SetGlossaryCreateFunc(handler func(string, string, string, [][2]string)) {
	ui.glossariesPage.SetGlossaryCreateFunc(handler)
}

func (ui *UI) SetGlossaryUpdateFunc(handler func(string, [][2]string)) {
	ui.glossariesPage.SetGlossaryUpdateFunc(handler)
}

func (ui *UI) SetGlossaryDeleteFunc(handler func(string)) {
	ui.glossariesPage.SetGlossaryDeleteFunc(handler)
}

// SetInputTextChangedFunc sets a handler that is called when the input text
// changes.
func (ui *UI) SetInputTextChangedFunc(handler func()) {
	ui.translatePage.inputTextArea.SetChangedFunc(handler)
}

func (ui *UI) GetInputText() string {
	return ui.translatePage.inputTextArea.GetText()
}

func (ui *UI) WriteOutputText(r io.Reader) error {
	w := ui.translatePage.outputTextView.BatchWriter()
	defer w.Close()
	_, err := io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

func (ui *UI) ClearOutputText() {
	ui.translatePage.outputTextView.Clear()
}

// Returns a new primitive which puts the provided one at the given position
// and sets its size to the given width and height.
func modal(p tview.Primitive, x, y, width, height int, focus bool) tview.Primitive {
	m := tview.NewGrid().
		SetColumns(x, width, 0).
		SetRows(y, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, focus)

	// m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	//     px, py, pw, ph := p.GetRect()
	//     ex, ey := event.Position()
	//     if px <= ex && ex <= px + pw && py <= ey && ey <= py + ph {
	//         consumed, capture = p.MouseHandler()(action, event, setFocus)
	//         return
	//     }
	//     return true, nil
	// })

	return m
}
