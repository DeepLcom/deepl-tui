package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/cluttrdev/deepl-go/deepl"
)

// GlossariesDialog is used to let the user choose a glossary displays the
// corresponding entries.
type GlossariesDialog struct {
	tview.Flex

	options [][2]string // list of id,name pairs

	dropDown *tview.DropDown
	table    *tview.Table
	buttons  *tview.Flex

	data     func(string) (deepl.GlossaryInfo, []deepl.GlossaryEntry)
	accepted func(string, string)
	cancel   func()
}

func newGlossariesDialog() *GlossariesDialog {
	w := &GlossariesDialog{
		Flex: *tview.NewFlex(),

		dropDown: tview.NewDropDown(),
		table:    tview.NewTable(),
	}

	// dropdown
	w.dropDown.SetLabel("Select: ")

	// table
	w.table.SetBorders(true)

	// buttons
	selectButton := tview.NewButton("Accept").
		SetSelectedFunc(func() {
			if w.accepted != nil {
				index, _ := w.dropDown.GetCurrentOption()
				var id, name string
				if index > 0 {
					id = w.options[index-1][0]
					name = w.options[index-1][1]
				}
				w.accepted(id, name)
			}
		})
	cancelButton := tview.NewButton("Cancel").
		SetSelectedFunc(func() {
			if w.cancel != nil {
				w.cancel()
			}
		})
	w.buttons = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(selectButton, 0, 1, true).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(cancelButton, 0, 1, false)

	// layout
	w.Flex.SetDirection(tview.FlexRow).
		AddItem(w.dropDown, 1, 0, true).
		AddItem(w.table, 0, 1, false).
		AddItem(w.buttons, 1, 0, false)

	return w
}

// SetOptions replaces all current glossary options with the ones provided.
func (w *GlossariesDialog) SetOptions(options [][2]string) *GlossariesDialog {
	w.options = options

	opts := make([]string, 1, len(options)+1)
	opts[0] = "        "
	for _, o := range w.options {
		opts = append(opts, o[1])
	}
	w.dropDown.SetOptions(opts, w.selectedFunc)
	return w
}

// SetDataFunc sets the handler that is used to get the glossary meta
// information and entries to display when a glossary is selected.
func (w *GlossariesDialog) SetDataFunc(data func(string) (deepl.GlossaryInfo, []deepl.GlossaryEntry)) *GlossariesDialog {
	w.data = data
	return w
}

// SetAcceptedFunc sets the handler which is called when the user accepts the
// current option by selecting the `accept` button.
func (w *GlossariesDialog) SetAcceptedFunc(accepted func(string, string)) *GlossariesDialog {
	w.accepted = accepted
	return w
}

// SetCancelFunc sets the handler which is called when the user selects the
// `cancel` button.
func (w *GlossariesDialog) SetCancelFunc(cancel func()) *GlossariesDialog {
	w.cancel = cancel
	return w
}

func (w *GlossariesDialog) selectedFunc(text string, index int) {
	w.table.Clear()
	if index > 0 {
		id := w.options[index-1][0]
		info, entries := w.data(id)
		if info.GlossaryId != "" {
			w.table.SetCell(0, 0, tview.NewTableCell(strings.ToUpper(info.SourceLang)).SetTextColor(tview.Styles.SecondaryTextColor).SetExpansion(1))
			w.table.SetCell(0, 1, tview.NewTableCell(strings.ToUpper(info.TargetLang)).SetTextColor(tview.Styles.SecondaryTextColor).SetExpansion(1))
		}
		for row, entry := range entries {
			w.table.SetCell(1+row, 0, tview.NewTableCell(entry.Source).SetExpansion(1))
			w.table.SetCell(1+row, 1, tview.NewTableCell(entry.Target).SetExpansion(1))
		}
	}
	w.table.ScrollToBeginning()
}

// GlossaryInfoForm displays glossary meta information in a form layout.
type GlossaryInfoForm struct {
	*tview.Form

	nameItem         *tview.InputField
	idItem           *tview.InputField
	sourceLangItem   *tview.DropDown
	targetLangItem   *tview.DropDown
	creationTimeItem *tview.InputField
	entryCountItem   *tview.InputField
}

func newGlossaryInfoForm() *GlossaryInfoForm {
	w := &GlossaryInfoForm{
		Form: tview.NewForm(),

		nameItem:         tview.NewInputField(),
		idItem:           tview.NewInputField(),
		sourceLangItem:   tview.NewDropDown(),
		targetLangItem:   tview.NewDropDown(),
		creationTimeItem: tview.NewInputField(),
		entryCountItem:   tview.NewInputField(),
	}

	w.nameItem.SetLabel("Name").SetFieldWidth(48)
	w.idItem.SetLabel("ID").SetDisabled(true)
	w.sourceLangItem.SetLabel("Source Lang").SetDisabled(false)
	w.targetLangItem.SetLabel("Target Lang").SetDisabled(false)
	w.creationTimeItem.SetLabel("Creation Time").SetDisabled(true)
	w.entryCountItem.SetLabel("Entry Count").SetDisabled(true)

	w.Form.AddFormItem(w.nameItem)
	w.Form.AddFormItem(w.idItem)
	w.Form.AddFormItem(w.sourceLangItem)
	w.Form.AddFormItem(w.targetLangItem)
	w.Form.AddFormItem(w.creationTimeItem)
	w.Form.AddFormItem(w.entryCountItem)

	return w
}

// SetLanguageOptions replaces all glossary source and target language options
// with the ones provided.
func (w *GlossaryInfoForm) SetLanguageOptions(langs []string) *GlossaryInfoForm {
	w.sourceLangItem.SetOptions(langs, nil)
	w.targetLangItem.SetOptions(langs, nil)
	return w
}

// SetInfo sets the glossary meta information.
func (w *GlossaryInfoForm) SetInfo(info deepl.GlossaryInfo) *GlossaryInfoForm {
	w.nameItem.SetText(info.Name)
	w.idItem.SetText(info.GlossaryId)
	w.sourceLangItem.
		SetTextOptions("", "", "", "", info.SourceLang).
		SetCurrentOption(-1)
	w.targetLangItem.
		SetTextOptions("", "", "", "", info.TargetLang).
		SetCurrentOption(-1)
	w.creationTimeItem.SetText(info.CreationTime)
	w.entryCountItem.SetText(fmt.Sprintf("%d", info.EntryCount))

	return w
}

// GlossaryEntryForm is used to manage glossary entries.
type GlossaryEntryForm struct {
	*tview.Form

	sourceItem *tview.InputField
	targetItem *tview.InputField
}

func newGlossaryEntryForm() *GlossaryEntryForm {
	w := &GlossaryEntryForm{
		Form: tview.NewForm(),

		sourceItem: tview.NewInputField(),
		targetItem: tview.NewInputField(),
	}

	w.sourceItem.SetLabel("Source").SetFieldWidth(64)
	w.targetItem.SetLabel("Target").SetFieldWidth(64)

	w.Form.AddFormItem(w.sourceItem)
	w.Form.AddFormItem(w.targetItem)

	return w
}

func (w *GlossaryEntryForm) SetEntry(entry deepl.GlossaryEntry) *GlossaryEntryForm {
	w.sourceItem.SetText(entry.Source)
	w.targetItem.SetText(entry.Target)
	return w
}

// GlossariesPage provides widgets to manage glossaries.
type GlossariesPage struct {
	*tview.Flex
	// *tview.Box

	list      *tview.List
	infoForm  *GlossaryInfoForm
	entryForm *GlossaryEntryForm
	table     *tview.Table

	data func(string) (deepl.GlossaryInfo, []deepl.GlossaryEntry)

	create func(name string, source string, target string, entries [][2]string)
	update func(id string, entries [][2]string)
	delete func(id string)
}

func newGlossariesPage(ui *UI) *GlossariesPage {
	w := &GlossariesPage{
		Flex: tview.NewFlex(),
		// Box: tview.NewBox(),

		list:      tview.NewList(),
		infoForm:  newGlossaryInfoForm(),
		entryForm: newGlossaryEntryForm(),
		table:     tview.NewTable(),
	}

	w.list.
		SetSelectedFocusOnly(false).
		SetChangedFunc(func(index int, name string, id string, _ rune) {
			if index == 0 && w.list.GetItemCount() > 0 {
				w.list.SetCurrentItem(1)
			}
		}).
		SetSelectedFunc(func(index int, name string, id string, _ rune) {
			w.selectedFunc(id, index)
		}).
		ShowSecondaryText(false)
	w.infoForm.
		AddButton("Create", w.onCreateGlossary).
		AddButton("Update", w.onUpdateGlossary).
		AddButton("Delete", w.onDeleteGlossary).
		SetHorizontal(false).
		SetItemPadding(0).
		SetTitle("Glossary Info").SetBorder(true)
	w.entryForm.
		AddButton("Create", w.onCreateEntry).
		AddButton("Update", w.onUpdateEntry).
		AddButton("Delete", w.onDeleteEntry).
		SetHorizontal(false).
		SetItemPadding(0)
	w.table.
		SetSelectable(true, false).
		SetSelectionChangedFunc(func(row, column int) {
			source := w.table.GetCell(row, 0).Text
			target := w.table.GetCell(row, 1).Text
			w.entryForm.sourceItem.SetText(source)
			w.entryForm.targetItem.SetText(target)
		}).
		SetSelectedStyle(tcell.StyleDefault.Foreground(tview.Styles.SecondaryTextColor).Bold(true)).
		SetBorders(true)

	newButton := tview.NewButton("New").
		SetSelectedFunc(func() {
			// select dummy list item 0 to signify preparation of a new glossary
			w.selectedFunc("", 0)
			ui.SetFocus(w.infoForm.nameItem)
		})

	leftLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(newButton, 1, 0, false).
		AddItem(w.list, 0, 1, false)
	leftLayout.SetTitle("Glossary List").SetBorder(true)

	entriesLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(w.entryForm, 6, 0, false).
		AddItem(w.table, 0, 1, false)
	entriesLayout.SetTitle("Glossary Entries").SetBorder(true)

	rightLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(w.infoForm, 12, 0, false).
		AddItem(entriesLayout, 0, 1, false)

	layout := tview.NewFlex()
	layout.SetDirection(tview.FlexColumn).
		AddItem(leftLayout, 32, 0, true).
		AddItem(rightLayout, 0, 1, false)
	w.Flex.AddItem(layout, 0, 1, true)

	w.registerKeyBindings(ui)

	return w
}

/*
func (w *GlossariesPage) Draw(screen tcell.Screen) {
	w.Box.DrawForSubclass(screen, w)
	x, y, width, height := w.GetInnerRect()

	w.hLayout.SetRect(x, y, width, height)
	w.hLayout.Draw(screen)
}

func (w *GlossariesPage) Focus(delegate func(p tview.Primitive)) {
    delegate(w.list)
}
*/

// SetOptions replaces the glossary list items with the ones provided.
func (w *GlossariesPage) SetOptions(options [][2]string) *GlossariesPage {
	var (
		id, name string
	)
	const nobind rune = 0
	w.list.Clear()
	for _, opt := range options {
		id = opt[0]
		name = opt[1]
		w.list.AddItem(name, id, nobind, nil)
	}
	w.list.InsertItem(0, "", "", 0, nil).SetCurrentItem(0)

	return w
}

// SetLanguageOptions replaces the glossary info form source and target language
// options with the ones provided.
func (w *GlossariesPage) SetLanguageOptions(langs []string) *GlossariesPage {
	w.infoForm.SetLanguageOptions(langs)
	return w
}

// SetGlossaryDataFunc sets the handler that is used to get the glossary meta
// information and entries to display when a glossary is selected.
func (w *GlossariesPage) SetGlossaryDataFunc(data func(string) (deepl.GlossaryInfo, []deepl.GlossaryEntry)) *GlossariesPage {
	w.data = data
	return w
}

// SetGlossaryCreateFunc sets the handler that is called when the user selects the
// glossary `create` button.
func (w *GlossariesPage) SetGlossaryCreateFunc(create func(string, string, string, [][2]string)) {
	w.create = create
}

// SetGlossaryUpdateFunc sets the handler that is called when the user selects the
// glossary `update` button.
func (w *GlossariesPage) SetGlossaryUpdateFunc(update func(string, [][2]string)) {
	w.update = update
}

// SetGlossaryDeleteFunc sets the handler that is called when the user selects the
// glossary `delete` button.
func (w *GlossariesPage) SetGlossaryDeleteFunc(del func(string)) {
	w.delete = del
}

func (w *GlossariesPage) registerKeyBindings(ui *UI) {
	w.Flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModAlt {
			switch event.Rune() {
			case 'e':
				ui.SetFocus(w.entryForm)
				return nil
			case 'i':
				ui.SetFocus(w.infoForm)
				return nil
			case 'l':
				ui.SetFocus(w.list)
				return nil
			case 't':
				ui.SetFocus(w.table)
				return nil
			}
		}
		return event
	})

}

func (w *GlossariesPage) selectedFunc(id string, index int) {
	var isIndex0 bool = (index == 0)
	w.infoForm.sourceLangItem.SetDisabled(!isIndex0)
	w.infoForm.targetLangItem.SetDisabled(!isIndex0)
	w.infoForm.GetButton(w.infoForm.GetButtonIndex("Create")).SetDisabled(!isIndex0)
	w.infoForm.GetButton(w.infoForm.GetButtonIndex("Update")).SetDisabled(isIndex0)
	w.infoForm.GetButton(w.infoForm.GetButtonIndex("Delete")).SetDisabled(isIndex0)

	w.table.Clear()
	if w.data != nil {
		info, entries := w.data(id)

		w.infoForm.SetInfo(info)
		for row, entry := range entries {
			w.table.SetCell(row, 0, tview.NewTableCell(entry.Source).SetExpansion(1))
			w.table.SetCell(row, 1, tview.NewTableCell(entry.Target).SetExpansion(1))
		}
		w.table.Select(w.table.GetRowCount(), 0) // `unselect`
	}
	w.table.ScrollToBeginning()
}

func (w *GlossariesPage) getTableEntries() [][2]string {
	entries := make([][2]string, 0, w.table.GetRowCount())
	for r := 0; r < w.table.GetRowCount(); r++ {
		entries = append(entries, [2]string{
			w.table.GetCell(r, 0).Text,
			w.table.GetCell(r, 1).Text,
		})
	}
	return entries
}

func (w *GlossariesPage) onCreateGlossary() {
	if w.create != nil {
		name := w.infoForm.nameItem.GetText()
		_, source := w.infoForm.sourceLangItem.GetCurrentOption()
		_, target := w.infoForm.targetLangItem.GetCurrentOption()
		entries := w.getTableEntries()

		w.create(name, source, target, entries)
	}
}

func (w *GlossariesPage) onUpdateGlossary() {
	if w.update != nil {
		id := w.infoForm.idItem.GetText()
		entries := w.getTableEntries()
		w.update(id, entries)
	}
}

func (w *GlossariesPage) onDeleteGlossary() {
	if w.delete != nil {
		id := w.infoForm.idItem.GetText()
		w.delete(id)
	}
}

func (w *GlossariesPage) onCreateEntry() {
	source := w.entryForm.sourceItem.GetText()
	target := w.entryForm.targetItem.GetText()

	w.table.
		InsertRow(0).
		SetCell(0, 0, tview.NewTableCell(source).SetExpansion(1)).
		SetCell(0, 1, tview.NewTableCell(target).SetExpansion(1)).
		Select(0, 0)
}

func (w *GlossariesPage) onUpdateEntry() {
	row, _ := w.table.GetSelection()

	source := w.entryForm.sourceItem.GetText()
	target := w.entryForm.targetItem.GetText()

	w.table.
		SetCell(row, 0, tview.NewTableCell(source).SetExpansion(1)).
		SetCell(row, 1, tview.NewTableCell(target).SetExpansion(1))
}

func (w *GlossariesPage) onDeleteEntry() {
	row, _ := w.table.GetSelection()

	w.table.RemoveRow(row)
}
