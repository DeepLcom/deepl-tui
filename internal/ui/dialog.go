package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Dialog struct {
	tview.Flex

	buttons *tview.Flex
	cancel  func()
}

func NewDialog(p tview.Primitive) *Dialog {
	dialog := &Dialog{
		Flex: *tview.NewFlex(),
	}
	dialog.Flex.SetDirection(tview.FlexRow)

	dialog.buttons = tview.NewFlex()

	dialog.Flex.
		AddItem(p, 0, 1, true).
		AddItem(dialog.buttons, 1, 0, false)

	dialog.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc && dialog.cancel != nil {
			dialog.cancel()
			return nil
		}
		return event
	})

	return dialog
}

func (d *Dialog) AddButton(label string, selected func()) *Dialog {
	button := tview.NewButton(label).
		SetSelectedFunc(selected)
	if d.buttons.GetItemCount() > 0 {
		d.buttons.AddItem(tview.NewBox(), 1, 0, false)
	}
	d.buttons.AddItem(button, 0, 1, false)
	return d
}

func (d *Dialog) SetCancelFunc(cancel func()) *Dialog {
	d.cancel = cancel
	return d
}
