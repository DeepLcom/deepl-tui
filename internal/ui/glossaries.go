package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/cluttrdev/deepl-go/deepl"
)

type GlossariesWidget struct {
	tview.Flex

	dropDown *tview.DropDown
	selected func(string, int)

	table *tview.Table
	data  func(string, int) (*deepl.GlossaryInfo, []deepl.GlossaryEntry)
}

type GlossaryTableContent struct {
	tview.TableContentReadOnly

	Info    *deepl.GlossaryInfo
	Entries []deepl.GlossaryEntry
}

func (td *GlossaryTableContent) GetCell(row, column int) *tview.TableCell {
	if row < 0 || row > len(td.Entries) || column < 0 || column >= 2 {
		return nil
	}

	cell := tview.NewTableCell("")

	if row == 0 {
		text := td.Info.SourceLang
		if column == 1 {
			text = td.Info.TargetLang
		}
		cell.
			SetText(strings.ToUpper(text)).
			SetStyle(tcell.StyleDefault.Foreground(tview.Styles.SecondaryTextColor))
	} else if column == 0 {
		cell.SetText(td.Entries[row-1].Source)
	} else if column == 1 {
		cell.SetText(td.Entries[row-1].Target)
	}

	return cell.SetExpansion(1)
}

func (td *GlossaryTableContent) GetRowCount() int {
	if td == nil {
		return 0
	}
	return len(td.Entries) + 1
}

func (td *GlossaryTableContent) GetColumnCount() int {
	return 2
}

func newGlossariesWidget(ui *UI) *GlossariesWidget {
	w := &GlossariesWidget{
		Flex: *tview.NewFlex(),
	}
	// layout
	w.Flex.SetDirection(tview.FlexRow)

	// dropdown
	w.dropDown = tview.NewDropDown().
		SetLabel("Select: ")
	w.Flex.AddItem(w.dropDown, 1, 0, true)

	// table
	w.table = tview.NewTable().
		SetFixed(1, 2).
		SetBorders(true)
	w.Flex.AddItem(w.table, 0, 1, true)

	return w
}

func (w *GlossariesWidget) GetCurrentOption() (int, string) {
	return w.dropDown.GetCurrentOption()
}

func (w *GlossariesWidget) SetOptions(options []string) *GlossariesWidget {
	opts := make([]string, 1, len(options)+1)
	opts[0] = "        "
	opts = append(opts, options...)
	w.dropDown.SetOptions(opts, w.selectedFunc)
	return w
}

func (w *GlossariesWidget) SetDataFunc(data func(string, int) (*deepl.GlossaryInfo, []deepl.GlossaryEntry)) *GlossariesWidget {
	w.data = data
	return w
}

func (w *GlossariesWidget) SetSelectedFunc(selected func(string, int)) *GlossariesWidget {
	w.selected = selected
	return w
}

func (w *GlossariesWidget) selectedFunc(text string, index int) {
	var content *GlossaryTableContent = nil
	if index > 0 {
		info, entries := w.data(text, index)
		if info != nil {
			content = &GlossaryTableContent{
				Info:    info,
				Entries: entries,
			}
		}
	}
	w.table.
		SetContent(content).
		ScrollToBeginning()

	if w.selected != nil {
		w.selected(text, index)
	}
}
